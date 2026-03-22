package approval

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	items []Approval
}

func (f fakeService) Create(context.Context, CreateApprovalCmd) (Approval, error) {
	return Approval{}, nil
}

func (f fakeService) Get(context.Context, string) (Approval, error) {
	return Approval{}, nil
}

func (f fakeService) ListByMission(context.Context, string) ([]Approval, error) {
	return f.items, nil
}

func TestListApprovalsByMission(t *testing.T) {
	handler := NewHandler(fakeService{
		items: []Approval{
			{ID: "approval_1", MissionID: "mission_1", Action: "repo_candidate_publish", Status: "pending"},
		},
	})

	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/missions/mission_1/approvals", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
