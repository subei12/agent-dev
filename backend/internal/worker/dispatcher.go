package worker

import "context"

type Dispatcher struct {
	runService *RunService
}

// NewDispatcher 创建并返回对应的组件。
func NewDispatcher(runService *RunService) *Dispatcher {
	return &Dispatcher{runService: runService}
}

// DispatchClaim 将当前工作负载派发给下一个执行器。
func (d *Dispatcher) DispatchClaim(ctx context.Context, claimID string) error {
	return d.runService.ExecuteClaimedTask(ctx, claimID)
}
