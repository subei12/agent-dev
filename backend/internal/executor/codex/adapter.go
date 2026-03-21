package codex

import "github.com/your-org/agent-platform/internal/executor"

type Adapter struct {
	executor.BaseAdapter
}

// NewAdapter 创建并返回对应的组件。
func NewAdapter() Adapter {
	return Adapter{BaseAdapter: executor.NewBaseAdapter("codex_cli")}
}
