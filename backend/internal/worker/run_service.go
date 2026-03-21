package worker

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type ClaimExecutionContext struct {
	ClaimID           string
	MissionID         string
	TaskItemID        string
	AgentID           string
	ExecutorProfileID string
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

type Store interface {
	GetClaimExecutionContext(context.Context, string) (ClaimExecutionContext, error)
	CreateRun(context.Context, string, string) (Run, error)
	CreateNodeRun(context.Context, string, string) (NodeRun, error)
	UpdateRunStatus(context.Context, string, string) error
}

type Launcher interface {
	RunTask(context.Context, ExecutorTaskInput) error
}

type RunService struct {
	store    Store
	launcher Launcher
}

func NewRunService(store Store, launcher Launcher) *RunService {
	return &RunService{store: store, launcher: launcher}
}

func (s *RunService) ExecuteClaimedTask(ctx context.Context, claimID string) error {
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

	if err := s.launcher.RunTask(ctx, ExecutorTaskInput{
		MissionID:         execCtx.MissionID,
		TaskItemID:        execCtx.TaskItemID,
		AgentID:           execCtx.AgentID,
		ExecutorProfileID: execCtx.ExecutorProfileID,
		Command:           "bash",
		Args:              []string{"-lc", "printf 'worker run\n'"},
	}); err != nil {
		_ = s.store.UpdateRunStatus(ctx, run.ID, "failed")
		return err
	}

	return s.store.UpdateRunStatus(ctx, run.ID, "succeeded")
}

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{queries: sqlc.New(pool)}
}

func (r *Repository) GetClaimExecutionContext(ctx context.Context, claimID string) (ClaimExecutionContext, error) {
	row, err := r.queries.GetTaskClaimExecutionContext(ctx, claimID)
	if err != nil {
		return ClaimExecutionContext{}, err
	}
	return ClaimExecutionContext{
		ClaimID:           row.ClaimID,
		MissionID:         row.MissionID,
		TaskItemID:        row.TaskItemID,
		AgentID:           row.AgentID,
		ExecutorProfileID: row.ExecutorProfileID,
	}, nil
}

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

func (r *Repository) UpdateRunStatus(ctx context.Context, runID, status string) error {
	_, err := r.queries.UpdateRunStatus(ctx, sqlc.UpdateRunStatusParams{
		ID:     runID,
		Status: status,
	})
	return err
}
