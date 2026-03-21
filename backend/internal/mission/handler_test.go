package mission

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	createMission func(context.Context, CreateMissionCmd) (Mission, error)
}

func (f fakeService) CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error) {
	return f.createMission(ctx, cmd)
}

func (f fakeService) GetMission(context.Context, string, string) (Mission, error) {
	return Mission{}, nil
}

func (f fakeService) Decide(context.Context, DecideMissionCmd) (MissionDecision, error) {
	return MissionDecision{}, nil
}

func TestCreateMission(t *testing.T) {
	svc := fakeService{
		createMission: func(_ context.Context, cmd CreateMissionCmd) (Mission, error) {
			return Mission{
				ID:        "mission_1",
				ProjectID: cmd.ProjectID,
				Title:     cmd.Title,
				Status:    "discussion",
			}, nil
		},
	}

	handler := NewHandler(svc)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	body := []byte(`{"title":"Build runtime board","teamId":"team_1"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/projects/proj_1/missions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}
