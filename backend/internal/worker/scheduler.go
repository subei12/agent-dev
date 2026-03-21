package worker

import (
	"context"
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

func NewScheduler(source ClaimSource, dispatcher ClaimDispatcher) *Scheduler {
	return &Scheduler{
		source:     source,
		dispatcher: dispatcher,
	}
}

func (s *Scheduler) PollOnce(ctx context.Context) error {
	claimIDs, err := s.source.ListActiveClaimIDs(ctx)
	if err != nil {
		return err
	}

	for _, claimID := range claimIDs {
		if err := s.dispatcher.DispatchClaim(ctx, claimID); err != nil {
			return err
		}
	}
	return nil
}

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
