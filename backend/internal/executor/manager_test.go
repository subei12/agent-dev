package executor

import (
	"context"
	"os/exec"
	"testing"

	"github.com/your-org/agent-platform/internal/runtime"
)

type fakeAdapter struct{}

func (fakeAdapter) Name() string { return "codex_cli" }

func (fakeAdapter) BuildCommand(_ context.Context, _ TaskExecutionInput) (*exec.Cmd, error) {
	return exec.Command("bash", "-lc", "printf 'assistant: done\n'"), nil
}

func (fakeAdapter) MapOutput(line []byte) ([]runtime.AppendEventCmd, []runtime.AppendTranscriptEntryCmd) {
	return []runtime.AppendEventCmd{
			{
				Level:    "info",
				Category: "model",
				Type:     "message_created",
				Title:    string(line),
			},
		}, []runtime.AppendTranscriptEntryCmd{
			{
				Role:        "assistant",
				EntryType:   "message",
				ContentText: string(line),
			},
		}
}

type fakeObserver struct {
	sessionStarted bool
	events         []runtime.AppendEventCmd
}

func (f *fakeObserver) StartSession(_ context.Context, cmd runtime.StartSessionCmd) (runtime.ExecutorSession, error) {
	f.sessionStarted = true
	return runtime.ExecutorSession{
		ID:        "session_1",
		MissionID: cmd.MissionID,
		AgentID:   cmd.AgentID,
	}, nil
}

func (f *fakeObserver) AppendEvent(_ context.Context, cmd runtime.AppendEventCmd) error {
	f.events = append(f.events, cmd)
	return nil
}

func (f *fakeObserver) AppendTranscriptEntry(context.Context, runtime.AppendTranscriptEntryCmd) error {
	return nil
}

func (f *fakeObserver) SealSession(context.Context, string) error { return nil }

type fakeRunner struct{}

func (fakeRunner) Run(_ context.Context, _ *exec.Cmd, onOutput func([]byte)) error {
	onOutput([]byte("assistant: done"))
	return nil
}

func TestManagerEmitsSessionStartedEvent(t *testing.T) {
	observer := &fakeObserver{}
	manager := NewManager(fakeAdapter{}, observer, fakeRunner{})

	err := manager.RunTask(context.Background(), TaskExecutionInput{
		MissionID:         "mission_1",
		TaskItemID:        "task_1",
		AgentID:           "agent_1",
		ExecutorProfileID: "exec_1",
		Command:           "bash",
		Args:              []string{"-lc", "printf 'assistant: done\n'"},
	})
	if err != nil {
		t.Fatalf("run task: %v", err)
	}

	if !observer.sessionStarted {
		t.Fatal("expected session to start")
	}

	var found bool
	for _, event := range observer.events {
		if event.Type == "session_started" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected session_started event")
	}
}
