package worker

import (
	"context"
	"testing"
)

type fakeClaimSource struct {
	claimIDs []string
}

// ListActiveClaimIDs 返回当前查询对应的集合结果。
func (f fakeClaimSource) ListActiveClaimIDs(context.Context) ([]string, error) {
	return f.claimIDs, nil
}

type fakeDispatcher struct {
	dispatched []string
	failClaims map[string]error
}

// DispatchClaim 将当前工作负载派发给下一个执行器。
func (f *fakeDispatcher) DispatchClaim(_ context.Context, claimID string) error {
	f.dispatched = append(f.dispatched, claimID)
	if f.failClaims != nil {
		if err, ok := f.failClaims[claimID]; ok {
			return err
		}
	}
	return nil
}

// TestSchedulerPollOnceDispatchesActiveClaims 验证该路径的预期行为。
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

// TestSchedulerPollOnceContinuesWhenSingleClaimFails 验证单条 claim 失败不会中断整轮调度。
func TestSchedulerPollOnceContinuesWhenSingleClaimFails(t *testing.T) {
	dispatcher := &fakeDispatcher{
		failClaims: map[string]error{
			"claim_1": context.DeadlineExceeded,
		},
	}
	scheduler := NewScheduler(fakeClaimSource{claimIDs: []string{"claim_1", "claim_2"}}, dispatcher)

	if err := scheduler.PollOnce(context.Background()); err != nil {
		t.Fatalf("poll once: %v", err)
	}

	if len(dispatcher.dispatched) != 2 {
		t.Fatalf("expected 2 dispatched claims, got %d", len(dispatcher.dispatched))
	}
	if dispatcher.dispatched[1] != "claim_2" {
		t.Fatalf("expected claim_2 to still be dispatched, got %v", dispatcher.dispatched)
	}
}
