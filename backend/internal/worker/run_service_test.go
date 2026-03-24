package worker

import (
	"context"
	"testing"
	"time"
)

type fakeStore struct {
	claimContext ClaimExecutionContext
	runs         []Run
	acquired     bool
	acquireOK    bool
	completed    bool
	failed       bool
	updatedTaskStatus map[string]string
	resolvedTasks     map[string]TaskTransitionTarget
	createdClaims     []ClaimTaskInput
	handoffs          []AdminHandoffInput
	missionDecisions  []MissionDecisionInput
	missionStatuses   map[string]string
}

// GetClaimExecutionContext 返回请求的资源或值。
func (f *fakeStore) GetClaimExecutionContext(context.Context, string) (ClaimExecutionContext, error) {
	return f.claimContext, nil
}

// AcquireClaimExecution 在可用时获取请求的执行锁。
func (f *fakeStore) AcquireClaimExecution(context.Context, string, string, time.Time) (bool, error) {
	f.acquired = true
	if !f.acquireOK {
		return false, nil
	}
	return true, nil
}

// TouchClaimHeartbeat 刷新当前心跳或时间戳。
func (f *fakeStore) TouchClaimHeartbeat(context.Context, string, string) error {
	return nil
}

// FailClaim 将当前工作流项标记为失败或可重试。
func (f *fakeStore) FailClaim(context.Context, string, string, int32) error {
	f.failed = true
	return nil
}

// CompleteClaim 将当前工作流项标记为已完成。
func (f *fakeStore) CompleteClaim(context.Context, string) error {
	f.completed = true
	return nil
}

// CreateRun 创建请求的资源或记录。
func (f *fakeStore) CreateRun(_ context.Context, missionID, taskItemID string) (Run, error) {
	run := Run{ID: "run_1", MissionID: missionID, TaskItemID: taskItemID, Status: "running"}
	f.runs = append(f.runs, run)
	return run, nil
}

// CreateNodeRun 创建请求的资源或记录。
func (f *fakeStore) CreateNodeRun(context.Context, string, string) (NodeRun, error) {
	return NodeRun{ID: "node_run_1", Status: "running"}, nil
}

// UpdateRunStatus 更新请求的资源状态。
func (f *fakeStore) UpdateRunStatus(context.Context, string, string) error { return nil }

// UpdateTaskStatus 更新请求的资源状态。
func (f *fakeStore) UpdateTaskStatus(_ context.Context, taskItemID, status string) error {
	if f.updatedTaskStatus == nil {
		f.updatedTaskStatus = map[string]string{}
	}
	f.updatedTaskStatus[taskItemID] = status
	return nil
}

// ResolveTaskByIdentifier 返回 Mission 内部任务引用对应的任务信息。
func (f *fakeStore) ResolveTaskByIdentifier(_ context.Context, _ string, identifier string) (TaskTransitionTarget, error) {
	return f.resolvedTasks[identifier], nil
}

// CreateTaskClaim 创建请求的资源或记录。
func (f *fakeStore) CreateTaskClaim(_ context.Context, input ClaimTaskInput) error {
	f.createdClaims = append(f.createdClaims, input)
	return nil
}

// CreateAdminHandoff 创建管理员回流记录。
func (f *fakeStore) CreateAdminHandoff(_ context.Context, input AdminHandoffInput) error {
	f.handoffs = append(f.handoffs, input)
	return nil
}

// CreateMissionDecision 创建 Mission 结项决策。
func (f *fakeStore) CreateMissionDecision(_ context.Context, input MissionDecisionInput) error {
	f.missionDecisions = append(f.missionDecisions, input)
	return nil
}

// UpdateMissionStatus 更新 Mission 状态。
func (f *fakeStore) UpdateMissionStatus(_ context.Context, missionID, status string) error {
	if f.missionStatuses == nil {
		f.missionStatuses = map[string]string{}
	}
	f.missionStatuses[missionID] = status
	return nil
}

type fakeLauncher struct {
	called bool
	err    error
}

// RunTask 执行当前组件的主循环或工作流。
func (f *fakeLauncher) RunTask(context.Context, ExecutorTaskInput) error {
	f.called = true
	return f.err
}

// TestClaimedTaskCreatesRunAndExecutorSession 验证该路径的预期行为。
func TestClaimedTaskCreatesRunAndExecutorSession(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:           "claim_1",
			MissionID:         "mission_1",
			TaskItemID:        "task_1",
			AgentID:           "agent_1",
			ExecutorProfileID: "exec_1",
		},
		acquireOK: true,
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}

	if len(store.runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(store.runs))
	}
	if !store.acquired {
		t.Fatal("expected claim execution lock to be acquired")
	}
	if !store.completed {
		t.Fatal("expected claim to be completed")
	}
	if !launcher.called {
		t.Fatal("expected executor launcher to be called")
	}
}

// TestClaimedTaskFailureMarksClaimForRetry 验证该路径的预期行为。
func TestClaimedTaskFailureMarksClaimForRetry(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:           "claim_1",
			MissionID:         "mission_1",
			TaskItemID:        "task_1",
			AgentID:           "agent_1",
			ExecutorProfileID: "exec_1",
		},
		acquireOK: true,
	}
	launcher := &fakeLauncher{err: context.DeadlineExceeded}
	svc := NewRunService(store, launcher)

	err := svc.ExecuteClaimedTask(context.Background(), "claim_1")
	if err == nil {
		t.Fatal("expected launcher error")
	}
	if !store.failed {
		t.Fatal("expected claim to be marked for retry/failure")
	}
	if store.completed {
		t.Fatal("did not expect completed claim on launcher failure")
	}
}

// TestClaimAlreadyLockedSkipsExecution 验证该路径的预期行为。
func TestClaimAlreadyLockedSkipsExecution(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:           "claim_1",
			MissionID:         "mission_1",
			TaskItemID:        "task_1",
			AgentID:           "agent_1",
			ExecutorProfileID: "exec_1",
		},
		acquireOK: false,
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}
	if launcher.called {
		t.Fatal("did not expect executor launch when lock is not acquired")
	}
}

// TestSuccessfulClaimAutoClaimsDownstreamTask 验证成功完成后会自动激活下游任务。
func TestSuccessfulClaimAutoClaimsDownstreamTask(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:            "claim_1",
			MissionID:          "mission_1",
			TaskItemID:         "task_design",
			TaskType:           "design",
			AgentID:            "agent_admin",
			ExecutorProfileID:  "exec_admin",
			DownstreamTaskRefs: []string{"开发实现与自测"},
			AdminAgentID:       "agent_admin",
		},
		acquireOK: true,
		resolvedTasks: map[string]TaskTransitionTarget{
			"开发实现与自测": {
				ID:              "task_code",
				AssignedAgentID: "agent_backend",
				Status:          "todo",
			},
		},
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}

	if got := store.updatedTaskStatus["task_design"]; got != "done" {
		t.Fatalf("expected task_design status done, got %q", got)
	}
	if len(store.createdClaims) != 1 {
		t.Fatalf("expected 1 downstream claim, got %d", len(store.createdClaims))
	}
	if store.createdClaims[0].TaskItemID != "task_code" {
		t.Fatalf("expected downstream task_code, got %q", store.createdClaims[0].TaskItemID)
	}
	if store.createdClaims[0].AgentID != "agent_backend" {
		t.Fatalf("expected downstream agent_backend, got %q", store.createdClaims[0].AgentID)
	}
	if got := store.updatedTaskStatus["task_code"]; got != "claimed" {
		t.Fatalf("expected downstream task claimed, got %q", got)
	}
	if len(store.missionDecisions) != 1 {
		t.Fatalf("expected 1 mission decision, got %d", len(store.missionDecisions))
	}
	if store.missionDecisions[0].Decision != "enter_implementation" {
		t.Fatalf("expected enter_implementation decision, got %q", store.missionDecisions[0].Decision)
	}
	if got := store.missionStatuses["mission_1"]; got != "implementation" {
		t.Fatalf("expected mission status implementation, got %q", got)
	}
}

// TestSuccessfulClaimWithoutDownstreamReturnsToAdmin 验证无下游时会自动回流管理员。
func TestSuccessfulClaimWithoutDownstreamReturnsToAdmin(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:           "claim_1",
			MissionID:         "mission_1",
			TaskItemID:        "task_test",
			AgentID:           "agent_backend",
			ExecutorProfileID: "exec_backend",
			AdminAgentID:      "agent_admin",
		},
		acquireOK: true,
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}

	if got := store.updatedTaskStatus["task_test"]; got != "handoff_pending" {
		t.Fatalf("expected task_test status handoff_pending, got %q", got)
	}
	if len(store.handoffs) != 1 {
		t.Fatalf("expected 1 admin handoff, got %d", len(store.handoffs))
	}
	if store.handoffs[0].AdminAgentID != "agent_admin" {
		t.Fatalf("expected admin handoff to agent_admin, got %q", store.handoffs[0].AdminAgentID)
	}
}

// TestSuccessfulAdminReviewCompletesMission 验证管理员最终复核任务会自动结项 Mission。
func TestSuccessfulAdminReviewCompletesMission(t *testing.T) {
	store := &fakeStore{
		claimContext: ClaimExecutionContext{
			ClaimID:           "claim_1",
			MissionID:         "mission_1",
			TaskItemID:        "task_review",
			TaskTitle:         "管理员复核与结项",
			AgentID:           "agent_admin",
			ExecutorProfileID: "exec_admin",
			AdminAgentID:      "agent_admin",
		},
		acquireOK: true,
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}

	if got := store.updatedTaskStatus["task_review"]; got != "done" {
		t.Fatalf("expected task_review status done, got %q", got)
	}
	if len(store.missionDecisions) != 1 {
		t.Fatalf("expected 1 mission decision, got %d", len(store.missionDecisions))
	}
	if store.missionDecisions[0].Decision != "complete_mission" {
		t.Fatalf("expected complete_mission decision, got %q", store.missionDecisions[0].Decision)
	}
	if got := store.missionStatuses["mission_1"]; got != "completed" {
		t.Fatalf("expected mission status completed, got %q", got)
	}
}
