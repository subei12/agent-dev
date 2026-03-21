package gemini

import "github.com/your-org/agent-platform/internal/executor"

type Adapter struct {
	executor.BaseAdapter
}

func NewAdapter() Adapter {
	return Adapter{BaseAdapter: executor.NewBaseAdapter("gemini_cli")}
}
