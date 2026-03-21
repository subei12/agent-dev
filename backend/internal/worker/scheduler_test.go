package worker

import (
	"context"
	"testing"
)

type fakeClaimSource struct {
	claimIDs []string
}

func (f fakeClaimSource) ListActiveClaimIDs(context.Context) ([]string, error) {
	return f.claimIDs, nil
}

type fakeDispatcher struct {
	dispatched []string
}

func (f *fakeDispatcher) DispatchClaim(_ context.Context, claimID string) error {
	f.dispatched = append(f.dispatched, claimID)
	return nil
}

func TestSchedulerPollOnceDispatchesActiveClaims(t *testing.T) {
	dispatcher := &fakeDispatcher{}
	scheduler := NewScheduler(fakeClaimSource{claimIDs: []string{"claim_1", "claim_2"}}, dispatcher)

	if err := scheduler.PollOnce(context.Background()); err != nil {
		t.Fatalf("poll once: %v", err)
	}

	if len(dispatcher.dispatched) != 2 {
		t.Fatalf("expected 2 dispatched claims, got %d", len(dispatcher.dispatched))
	}
}
