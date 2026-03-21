package worker

import (
	"context"

	executorpkg "github.com/your-org/agent-platform/internal/executor"
)

// ExecutorLauncher adapts worker task execution requests to the executor package contract.
type ExecutorLauncher struct {
	manager *executorpkg.Manager
}

func NewExecutorLauncher(manager *executorpkg.Manager) *ExecutorLauncher {
	return &ExecutorLauncher{manager: manager}
}

func (l *ExecutorLauncher) RunTask(ctx context.Context, input ExecutorTaskInput) error {
	return l.manager.RunTask(ctx, executorpkg.TaskExecutionInput{
		MissionID:         input.MissionID,
		TaskItemID:        input.TaskItemID,
		AgentID:           input.AgentID,
		ExecutorProfileID: input.ExecutorProfileID,
		Command:           input.Command,
		Args:              input.Args,
	})
}
