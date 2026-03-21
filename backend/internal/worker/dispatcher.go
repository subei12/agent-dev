package worker

import "context"

type Dispatcher struct {
	runService *RunService
}

func NewDispatcher(runService *RunService) *Dispatcher {
	return &Dispatcher{runService: runService}
}

func (d *Dispatcher) DispatchClaim(ctx context.Context, claimID string) error {
	return d.runService.ExecuteClaimedTask(ctx, claimID)
}
