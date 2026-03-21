package discussion

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	createSession func(context.Context, CreateSessionCmd) (Session, error)
}

func (f fakeService) CreateSession(ctx context.Context, cmd CreateSessionCmd) (Session, error) {
	return f.createSession(ctx, cmd)
}

func (f fakeService) ListSessions(context.Context, string) ([]Session, error) {
	return nil, nil
}

func (f fakeService) CreateRound(context.Context, CreateRoundCmd) (Round, error) {
	return Round{}, nil
}

func TestCreateDiscussionSession(t *testing.T) {
	svc := fakeService{
		createSession: func(_ context.Context, cmd CreateSessionCmd) (Session, error) {
			return Session{
				ID:        "session_1",
				MissionID: cmd.MissionID,
				Topic:     cmd.Topic,
				Status:    "open",
			}, nil
		},
	}

	handler := NewHandler(svc)
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	body := []byte(`{"topic":"需求澄清","initiatedByAgentId":"agent_admin"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/projects/proj_1/missions/mission_1/discussions", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}
