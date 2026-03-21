package executor

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/your-org/agent-platform/internal/runtime"
)

type RuntimeObserver interface {
	StartSession(context.Context, runtime.StartSessionCmd) (runtime.ExecutorSession, error)
	AppendEvent(context.Context, runtime.AppendEventCmd) error
	AppendTranscriptEntry(context.Context, runtime.AppendTranscriptEntryCmd) error
	SealSession(context.Context, string) error
}

type Runner interface {
	Run(context.Context, *exec.Cmd, func([]byte)) error
}

type Manager struct {
	adapter  Adapter
	observer RuntimeObserver
	runner   Runner
}

// NewManager 创建并返回对应的组件。
func NewManager(adapter Adapter, observer RuntimeObserver, runner Runner) *Manager {
	return &Manager{
		adapter:  adapter,
		observer: observer,
		runner:   runner,
	}
}

// RunTask 执行当前组件的主循环或工作流。
func (m *Manager) RunTask(ctx context.Context, input TaskExecutionInput) error {
	session, err := m.observer.StartSession(ctx, runtime.StartSessionCmd{
		MissionID:         input.MissionID,
		TaskItemID:        input.TaskItemID,
		AgentID:           input.AgentID,
		ExecutorProfileID: input.ExecutorProfileID,
		Backend:           m.adapter.Name(),
		Status:            "starting",
	})
	if err != nil {
		return err
	}

	if err := m.observer.AppendEvent(ctx, runtime.AppendEventCmd{
		MissionID:         input.MissionID,
		AgentID:           input.AgentID,
		ExecutorSessionID: session.ID,
		TaskItemID:        input.TaskItemID,
		Level:             "info",
		Category:          "session",
		Type:              "session_started",
		Title:             "Session started",
	}); err != nil {
		return err
	}

	cmd, err := m.adapter.BuildCommand(ctx, input)
	if err != nil {
		return err
	}

	if err := m.runner.Run(ctx, cmd, func(line []byte) {
		events, transcripts := m.adapter.MapOutput(line)
		for _, event := range events {
			event.MissionID = input.MissionID
			event.AgentID = input.AgentID
			event.ExecutorSessionID = session.ID
			event.TaskItemID = input.TaskItemID
			_ = m.observer.AppendEvent(ctx, event)
		}
		for _, entry := range transcripts {
			entry.ExecutorSessionID = session.ID
			_ = m.observer.AppendTranscriptEntry(ctx, entry)
		}
	}); err != nil {
		return fmt.Errorf("run task: %w", err)
	}

	return m.observer.SealSession(ctx, session.ID)
}
