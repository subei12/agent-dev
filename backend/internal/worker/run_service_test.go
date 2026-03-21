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
}

func (f *fakeStore) GetClaimExecutionContext(context.Context, string) (ClaimExecutionContext, error) {
	return f.claimContext, nil
}

func (f *fakeStore) AcquireClaimExecution(context.Context, string, string, time.Time) (bool, error) {
	f.acquired = true
	if !f.acquireOK {
		return false, nil
	}
	return true, nil
}

func (f *fakeStore) TouchClaimHeartbeat(context.Context, string, string) error {
	return nil
}

func (f *fakeStore) FailClaim(context.Context, string, string, int32) error {
	f.failed = true
	return nil
}

func (f *fakeStore) CompleteClaim(context.Context, string) error {
	f.completed = true
	return nil
}

func (f *fakeStore) CreateRun(_ context.Context, missionID, taskItemID string) (Run, error) {
	run := Run{ID: "run_1", MissionID: missionID, TaskItemID: taskItemID, Status: "running"}
	f.runs = append(f.runs, run)
	return run, nil
}

func (f *fakeStore) CreateNodeRun(context.Context, string, string) (NodeRun, error) {
	return NodeRun{ID: "node_run_1", Status: "running"}, nil
}

func (f *fakeStore) UpdateRunStatus(context.Context, string, string) error { return nil }

type fakeLauncher struct {
	called bool
	err    error
}

func (f *fakeLauncher) RunTask(context.Context, ExecutorTaskInput) error {
	f.called = true
	return f.err
}

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
