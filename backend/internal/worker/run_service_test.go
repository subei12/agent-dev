package worker

import (
	"context"
	"testing"
)

type fakeStore struct {
	claimContext ClaimExecutionContext
	runs         []Run
}

func (f *fakeStore) GetClaimExecutionContext(context.Context, string) (ClaimExecutionContext, error) {
	return f.claimContext, nil
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
}

func (f *fakeLauncher) RunTask(context.Context, ExecutorTaskInput) error {
	f.called = true
	return nil
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
	}
	launcher := &fakeLauncher{}
	svc := NewRunService(store, launcher)

	if err := svc.ExecuteClaimedTask(context.Background(), "claim_1"); err != nil {
		t.Fatalf("execute claimed task: %v", err)
	}

	if len(store.runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(store.runs))
	}
	if !launcher.called {
		t.Fatal("expected executor launcher to be called")
	}
}
