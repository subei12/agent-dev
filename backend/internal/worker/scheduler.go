package worker

import (
	"context"
	"log"
	"time"
)

type ClaimSource interface {
	ListActiveClaimIDs(context.Context) ([]string, error)
}

type ClaimDispatcher interface {
	DispatchClaim(context.Context, string) error
}

type Scheduler struct {
	source     ClaimSource
	dispatcher ClaimDispatcher
}

// NewScheduler 创建并返回对应的组件。
func NewScheduler(source ClaimSource, dispatcher ClaimDispatcher) *Scheduler {
	return &Scheduler{
		source:     source,
		dispatcher: dispatcher,
	}
}

// PollOnce 执行一次调度轮询。
func (s *Scheduler) PollOnce(ctx context.Context) error {
	claimIDs, err := s.source.ListActiveClaimIDs(ctx)
	if err != nil {
		return err
	}

	for _, claimID := range claimIDs {
		if err := s.dispatcher.DispatchClaim(ctx, claimID); err != nil {
			log.Printf("worker dispatch claim %s failed: %v", claimID, err)
			continue
		}
	}
	return nil
}

// Run 执行当前组件的主循环或工作流。
func (s *Scheduler) Run(ctx context.Context, interval time.Duration) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if err := s.PollOnce(ctx); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
