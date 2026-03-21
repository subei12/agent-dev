package executor

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/your-org/agent-platform/internal/runtime"
)

type TaskExecutionInput struct {
	MissionID         string
	TaskItemID        string
	AgentID           string
	ExecutorProfileID string
	Command           string
	Args              []string
}

type Adapter interface {
	Name() string
	BuildCommand(context.Context, TaskExecutionInput) (*exec.Cmd, error)
	MapOutput(line []byte) ([]runtime.AppendEventCmd, []runtime.AppendTranscriptEntryCmd)
}

type BaseAdapter struct {
	name string
}

func NewBaseAdapter(name string) BaseAdapter {
	return BaseAdapter{name: name}
}

func (a BaseAdapter) Name() string {
	return a.name
}

func (a BaseAdapter) BuildCommand(ctx context.Context, input TaskExecutionInput) (*exec.Cmd, error) {
	if input.Command == "" {
		return nil, fmt.Errorf("command is required")
	}
	return exec.CommandContext(ctx, input.Command, input.Args...), nil
}

func (a BaseAdapter) MapOutput(line []byte) ([]runtime.AppendEventCmd, []runtime.AppendTranscriptEntryCmd) {
	text := string(line)
	return []runtime.AppendEventCmd{
			{
				Level:    "info",
				Category: "model",
				Type:     "message_created",
				Title:    "Model output",
				Summary:  text,
			},
		}, []runtime.AppendTranscriptEntryCmd{
			{
				Role:        "assistant",
				EntryType:   "message",
				ContentText: text,
			},
		}
}
