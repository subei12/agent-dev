package runtime

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/your-org/agent-platform/internal/authz"
)

type fakeRuntimeService struct{}
type fakeAuthorizer struct{ err error }

func (fakeRuntimeService) StartSession(context.Context, StartSessionCmd) (ExecutorSession, error) {
	return ExecutorSession{}, nil
}
func (fakeRuntimeService) AppendEvent(context.Context, AppendEventCmd) error { return nil }
func (fakeRuntimeService) AppendTranscriptEntry(context.Context, AppendTranscriptEntryCmd) error {
	return nil
}
func (fakeRuntimeService) SealSession(context.Context, string) error { return nil }
func (fakeRuntimeService) ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error) {
	return nil, nil
}
func (fakeRuntimeService) GetSession(context.Context, string) (ExecutorSession, error) {
	return ExecutorSession{}, nil
}
func (fakeRuntimeService) ListSessionEvents(context.Context, string) ([]RuntimeEvent, error) {
	return nil, nil
}
func (fakeRuntimeService) GetTranscriptView(context.Context, string, string, string) (TranscriptView, error) {
	return TranscriptView{}, nil
}
func (fakeRuntimeService) ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error) {
	return nil, nil
}

func (f fakeAuthorizer) Require(context.Context, string, string, string) error {
	return f.err
}

func TestGetTranscriptRequiresScopeForRedactedView(t *testing.T) {
	handler := NewHandler(fakeRuntimeService{}, fakeAuthorizer{err: authz.ErrForbidden})
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/executor-sessions/session_1/transcript?view=redacted", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}

func TestGetTranscriptAllowsRedactedViewWhenAuthorized(t *testing.T) {
	handler := NewHandler(fakeRuntimeService{}, fakeAuthorizer{})
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/executor-sessions/session_1/transcript?view=redacted", nil)
	req.Header.Set("X-Actor-Id", "demo-user")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
