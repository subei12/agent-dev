package worker

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type ClaimExecutionContext struct {
	ClaimID            string
	MissionID          string
	TaskItemID         string
	TaskTitle          string
	AgentID            string
	AdminAgentID       string
	DownstreamTaskRefs []string
	ExecutorProfileID  string
	Command            string
	Args               []string
}

type Run struct {
	ID         string
	MissionID  string
	TaskItemID string
	Status     string
}

type NodeRun struct {
	ID     string
	RunID  string
	Status string
}

type ExecutorTaskInput struct {
	MissionID         string
	TaskItemID        string
	AgentID           string
	ExecutorProfileID string
	Command           string
	Args              []string
}

type TaskTransitionTarget struct {
	ID              string
	Title           string
	Status          string
	AssignedAgentID string
}

type ClaimTaskInput struct {
	TaskItemID   string
	AgentID      string
	ClaimReason  string
}

type AdminHandoffInput struct {
	MissionID     string
	TaskItemID    string
	FromAgentID   string
	AdminAgentID  string
	Summary       string
}

type MissionDecisionInput struct {
	MissionID        string
	DecidedByAgentID string
	Decision         string
	Summary          string
}

type Store interface {
	GetClaimExecutionContext(context.Context, string) (ClaimExecutionContext, error)
	AcquireClaimExecution(context.Context, string, string, time.Time) (bool, error)
	TouchClaimHeartbeat(context.Context, string, string) error
	FailClaim(context.Context, string, string, int32) error
	CompleteClaim(context.Context, string) error
	CreateRun(context.Context, string, string) (Run, error)
	CreateNodeRun(context.Context, string, string) (NodeRun, error)
	UpdateRunStatus(context.Context, string, string) error
	UpdateTaskStatus(context.Context, string, string) error
	ResolveTaskByIdentifier(context.Context, string, string) (TaskTransitionTarget, error)
	CreateTaskClaim(context.Context, ClaimTaskInput) error
	CreateAdminHandoff(context.Context, AdminHandoffInput) error
	CreateMissionDecision(context.Context, MissionDecisionInput) error
	UpdateMissionStatus(context.Context, string, string) error
}

type Launcher interface {
	RunTask(context.Context, ExecutorTaskInput) error
}

type RunService struct {
	store    Store
	launcher Launcher
}

const maxClaimExecutionAttempts int32 = 3
const claimHeartbeatInterval = 5 * time.Second
const staleClaimLockAfter = 30 * time.Second

// NewRunService 创建并返回对应的组件。
func NewRunService(store Store, launcher Launcher) *RunService {
	return &RunService{store: store, launcher: launcher}
}

// ExecuteClaimedTask 获取 claim 锁、执行任务，并在结束后更新 claim 状态。
func (s *RunService) ExecuteClaimedTask(ctx context.Context, claimID string) error {
	// 1. 先获取 claim 执行锁；如果已有其他 worker 持有该锁，则直接跳过。
	lockToken := uuid.NewString()
	acquired, err := s.store.AcquireClaimExecution(ctx, claimID, lockToken, time.Now().Add(-staleClaimLockAfter))
	if err != nil {
		return err
	}
	if !acquired {
		return nil
	}

	// 2. 在执行期间持续刷新 claim 心跳，避免被 stale reclaim 误判回收。
	heartbeatCtx, stopHeartbeat := context.WithCancel(ctx)
	defer stopHeartbeat()
	go s.runHeartbeatLoop(heartbeatCtx, claimID, lockToken)

	// 3. 创建 run/node_run 记录，并把实际执行委托给执行器。
	execCtx, err := s.store.GetClaimExecutionContext(ctx, claimID)
	if err != nil {
		return err
	}

	run, err := s.store.CreateRun(ctx, execCtx.MissionID, execCtx.TaskItemID)
	if err != nil {
		return err
	}
	if _, err := s.store.CreateNodeRun(ctx, run.ID, "task_execute"); err != nil {
		return err
	}

	// 4. 若执行失败则回退 claim 并记录错误；成功则更新任务状态并推进后续任务。
	if err := s.launcher.RunTask(ctx, ExecutorTaskInput{
		MissionID:         execCtx.MissionID,
		TaskItemID:        execCtx.TaskItemID,
		AgentID:           execCtx.AgentID,
		ExecutorProfileID: execCtx.ExecutorProfileID,
		Command:           execCtx.Command,
		Args:              execCtx.Args,
	}); err != nil {
		_ = s.store.UpdateRunStatus(ctx, run.ID, "failed")
		_ = s.store.FailClaim(ctx, claimID, err.Error(), maxClaimExecutionAttempts)
		return err
	}

	if err := s.store.UpdateRunStatus(ctx, run.ID, "succeeded"); err != nil {
		return err
	}
	if err := s.advanceTaskChain(ctx, execCtx); err != nil {
		return err
	}
	return s.store.CompleteClaim(ctx, claimID)
}

// advanceTaskChain 在当前任务成功后自动激活下游任务或回流管理员。
func (s *RunService) advanceTaskChain(ctx context.Context, execCtx ClaimExecutionContext) error {
	// 1. 有下游任务时，先完成当前任务，再为下游创建 claim。
	if len(execCtx.DownstreamTaskRefs) > 0 {
		if err := s.store.UpdateTaskStatus(ctx, execCtx.TaskItemID, "done"); err != nil {
			return err
		}

		for _, taskRef := range execCtx.DownstreamTaskRefs {
			target, err := s.store.ResolveTaskByIdentifier(ctx, execCtx.MissionID, taskRef)
			if err != nil {
				return err
			}
			if target.Status != "todo" {
				continue
			}
			if target.AssignedAgentID == "" {
				continue
			}
			if err := s.store.UpdateTaskStatus(ctx, target.ID, "claimed"); err != nil {
				return err
			}
			if err := s.store.CreateTaskClaim(ctx, ClaimTaskInput{
				TaskItemID:  target.ID,
				AgentID:     target.AssignedAgentID,
				ClaimReason: "auto progression",
			}); err != nil {
				return err
			}
		}
		return nil
	}

	// 2. 无下游且当前不是管理员任务时，自动回流管理员 Agent 等待裁决。
	if execCtx.AdminAgentID != "" && execCtx.AgentID != execCtx.AdminAgentID {
		if err := s.store.UpdateTaskStatus(ctx, execCtx.TaskItemID, "handoff_pending"); err != nil {
			return err
		}
		return s.store.CreateAdminHandoff(ctx, AdminHandoffInput{
			MissionID:    execCtx.MissionID,
			TaskItemID:   execCtx.TaskItemID,
			FromAgentID:  execCtx.AgentID,
			AdminAgentID: execCtx.AdminAgentID,
			Summary:      "当前步骤已完成，自动回流管理员 Agent 决策是否结项或追加检查。",
		})
	}

	// 3. 管理员任务无下游时，先标记为完成，后续由管理员决策链补结项。
	if err := s.store.UpdateTaskStatus(ctx, execCtx.TaskItemID, "done"); err != nil {
		return err
	}
	if err := s.store.CreateMissionDecision(ctx, MissionDecisionInput{
		MissionID:        execCtx.MissionID,
		DecidedByAgentID: execCtx.AdminAgentID,
		Decision:         "complete_mission",
		Summary:          "管理员 Agent 已完成最终复核，Mission 自动结项。",
	}); err != nil {
		return err
	}
	return s.store.UpdateMissionStatus(ctx, execCtx.MissionID, "completed")
}

// runHeartbeatLoop 在执行上下文结束前持续刷新 claim 心跳。
func (s *RunService) runHeartbeatLoop(ctx context.Context, claimID, lockToken string) {
	ticker := time.NewTicker(claimHeartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = s.store.TouchClaimHeartbeat(ctx, claimID, lockToken)
		}
	}
}

type Repository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

// NewRepository 创建并返回对应的组件。
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

// GetClaimExecutionContext 返回请求的资源或值。
func (r *Repository) GetClaimExecutionContext(ctx context.Context, claimID string) (ClaimExecutionContext, error) {
	row, err := r.queries.GetTaskClaimExecutionContext(ctx, claimID)
	if err != nil {
		return ClaimExecutionContext{}, err
	}
	return ClaimExecutionContext{
		ClaimID:            row.ClaimID,
		MissionID:          row.MissionID,
		TaskItemID:         row.TaskItemID,
		TaskTitle:          row.TaskTitle,
		AgentID:            row.AgentID,
		AdminAgentID:       row.AdminAgentID,
		DownstreamTaskRefs: jsonStringSlice(row.DownstreamTaskIdsJson),
		ExecutorProfileID:  row.ExecutorProfileID,
		Command:            row.Command,
		Args:               jsonStringSlice(row.ArgsJson),
	}, nil
}

// ListActiveClaimIDs 返回当前查询对应的集合结果。
func (r *Repository) ListActiveClaimIDs(ctx context.Context) ([]string, error) {
	return r.queries.ListActiveTaskClaimIDs(ctx)
}

// AcquireClaimExecution 在可用时获取请求的执行锁。
func (r *Repository) AcquireClaimExecution(ctx context.Context, claimID, lockToken string, staleBefore time.Time) (bool, error) {
	_, err := r.queries.AcquireTaskClaimExecution(ctx, sqlc.AcquireTaskClaimExecutionParams{
		ID:                 claimID,
		ExecutionLockToken: textValue(lockToken),
		LastHeartbeatAt:    timestamptzValue(staleBefore),
	})
	if err == nil {
		return true, nil
	}
	if err.Error() == "no rows in result set" {
		return false, nil
	}
	return false, err
}

// TouchClaimHeartbeat 刷新当前心跳或时间戳。
func (r *Repository) TouchClaimHeartbeat(ctx context.Context, claimID, lockToken string) error {
	_, err := r.queries.UpdateTaskClaimHeartbeat(ctx, sqlc.UpdateTaskClaimHeartbeatParams{
		ID:                 claimID,
		ExecutionLockToken: textValue(lockToken),
	})
	return err
}

// FailClaim 将当前工作流项标记为失败或可重试。
func (r *Repository) FailClaim(ctx context.Context, claimID, lastError string, maxAttempts int32) error {
	_, err := r.queries.FailTaskClaimExecution(ctx, sqlc.FailTaskClaimExecutionParams{
		ID:           claimID,
		AttemptCount: maxAttempts,
		LastError:    textValue(lastError),
	})
	return err
}

// CompleteClaim 将当前工作流项标记为已完成。
func (r *Repository) CompleteClaim(ctx context.Context, claimID string) error {
	_, err := r.queries.CompleteTaskClaim(ctx, claimID)
	return err
}

// CreateRun 创建请求的资源或记录。
func (r *Repository) CreateRun(ctx context.Context, missionID, taskItemID string) (Run, error) {
	row, err := r.queries.CreateRun(ctx, sqlc.CreateRunParams{
		ID:         uuid.NewString(),
		MissionID:  missionID,
		TaskItemID: textValue(taskItemID),
		Status:     "running",
		InputsJson: json.RawMessage("{}"),
	})
	if err != nil {
		return Run{}, err
	}
	return Run{ID: row.ID, MissionID: row.MissionID, TaskItemID: stringValue(row.TaskItemID), Status: row.Status}, nil
}

// CreateNodeRun 创建请求的资源或记录。
func (r *Repository) CreateNodeRun(ctx context.Context, runID, nodeKey string) (NodeRun, error) {
	row, err := r.queries.CreateNodeRun(ctx, sqlc.CreateNodeRunParams{
		ID:                uuid.NewString(),
		RunID:             runID,
		NodeKey:           nodeKey,
		Status:            "running",
		ExecutorSessionID: textValue(""),
	})
	if err != nil {
		return NodeRun{}, err
	}
	return NodeRun{ID: row.ID, RunID: row.RunID, Status: row.Status}, nil
}

// UpdateRunStatus 更新请求的资源状态。
func (r *Repository) UpdateRunStatus(ctx context.Context, runID, status string) error {
	_, err := r.queries.UpdateRunStatus(ctx, sqlc.UpdateRunStatusParams{
		ID:     runID,
		Status: status,
	})
	return err
}

// UpdateTaskStatus 更新请求的资源状态。
func (r *Repository) UpdateTaskStatus(ctx context.Context, taskItemID, status string) error {
	_, err := sqlc.New(r.pool).UpdateTaskItemStatus(ctx, sqlc.UpdateTaskItemStatusParams{
		ID:     taskItemID,
		Status: status,
	})
	return err
}

// ResolveTaskByIdentifier 返回 Mission 内的目标任务信息。
func (r *Repository) ResolveTaskByIdentifier(ctx context.Context, missionID, identifier string) (TaskTransitionTarget, error) {
	row, err := sqlc.New(r.pool).GetMissionTaskByIdentifier(ctx, sqlc.GetMissionTaskByIdentifierParams{
		MissionID: missionID,
		ID:        identifier,
	})
	if err != nil {
		return TaskTransitionTarget{}, err
	}
	return TaskTransitionTarget{
		ID:              row.ID,
		Title:           row.Title,
		Status:          row.Status,
		AssignedAgentID: stringValue(row.AssignedAgentID),
	}, nil
}

// CreateTaskClaim 创建自动推进产生的任务 claim。
func (r *Repository) CreateTaskClaim(ctx context.Context, input ClaimTaskInput) error {
	_, err := sqlc.New(r.pool).CreateTaskClaim(ctx, sqlc.CreateTaskClaimParams{
		ID:          uuid.NewString(),
		TaskItemID:  input.TaskItemID,
		AgentID:     input.AgentID,
		Status:      "active",
		ClaimReason: textValue(input.ClaimReason),
	})
	return err
}

// CreateAdminHandoff 在无下游时创建回流管理员的 handoff 记录。
func (r *Repository) CreateAdminHandoff(ctx context.Context, input AdminHandoffInput) error {
	_, err := sqlc.New(r.pool).CreateTaskHandoff(ctx, sqlc.CreateTaskHandoffParams{
		ID:                         uuid.NewString(),
		TaskItemID:                 input.TaskItemID,
		FromAgentID:                input.FromAgentID,
		ToAgentID:                  textValue(input.AdminAgentID),
		ToAdminAgent:               true,
		Summary:                    input.Summary,
		OutputDocumentVersionIdsJson: []byte("[]"),
		OutputRepoCandidateIdsJson:   []byte("[]"),
		OutputArtifactVersionIdsJson: []byte("[]"),
		ValidationSummary:            textValue(""),
		RiskSummary:                  textValue(""),
		RecommendedNextAction:        textValue("管理员 Agent 判断是否完成或追加检查"),
		Status:                       "pending",
	})
	return err
}

// CreateMissionDecision 写入管理员最终决策。
func (r *Repository) CreateMissionDecision(ctx context.Context, input MissionDecisionInput) error {
	_, err := r.queries.CreateMissionDecision(ctx, sqlc.CreateMissionDecisionParams{
		ID:                           uuid.NewString(),
		MissionID:                    input.MissionID,
		DecidedByAgentID:             input.DecidedByAgentID,
		Decision:                     input.Decision,
		Summary:                      input.Summary,
		RelatedTaskItemID:            textValue(""),
		RelatedDocumentVersionIdsJson: []byte("[]"),
		RelatedRepoCandidateIdsJson:   []byte("[]"),
	})
	return err
}

// UpdateMissionStatus 更新 Mission 主状态。
func (r *Repository) UpdateMissionStatus(ctx context.Context, missionID, status string) error {
	_, err := r.queries.UpdateMissionStatus(ctx, sqlc.UpdateMissionStatusParams{
		ID:     missionID,
		Status: status,
	})
	return err
}

// jsonStringSlice 实现当前函数行为。
func jsonStringSlice(raw []byte) []string {
	if len(raw) == 0 {
		return nil
	}
	var values []string
	if err := json.Unmarshal(raw, &values); err != nil {
		return nil
	}
	return values
}
