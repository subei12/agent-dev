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

// StartSession 实现当前函数行为。
func (fakeRuntimeService) StartSession(context.Context, StartSessionCmd) (ExecutorSession, error) {
	return ExecutorSession{}, nil
}
// AppendEvent 向当前流追加新的事件或转录条目。
func (fakeRuntimeService) AppendEvent(context.Context, AppendEventCmd) error { return nil }
// AppendTranscriptEntry 向当前流追加新的事件或转录条目。
func (fakeRuntimeService) AppendTranscriptEntry(context.Context, AppendTranscriptEntryCmd) error {
	return nil
}
// SealSession 完成当前运行时会话的收尾。
func (fakeRuntimeService) SealSession(context.Context, string) error { return nil }
// ListMissionRuntimes 返回当前查询对应的集合结果。
func (fakeRuntimeService) ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error) {
	return nil, nil
}
// GetSession 返回请求的资源或值。
func (fakeRuntimeService) GetSession(context.Context, string) (ExecutorSession, error) {
	return ExecutorSession{}, nil
}
// ListSessionEvents 返回当前查询对应的集合结果。
func (fakeRuntimeService) ListSessionEvents(context.Context, string) ([]RuntimeEvent, error) {
	return nil, nil
}
// GetTranscriptView 返回请求的资源或值。
func (fakeRuntimeService) GetTranscriptView(context.Context, string, string, string) (TranscriptView, error) {
	return TranscriptView{}, nil
}
// ListAccessAudits 返回当前查询对应的集合结果。
func (fakeRuntimeService) ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error) {
	return nil, nil
}

// Require 校验请求的授权能力。
func (f fakeAuthorizer) Require(context.Context, string, string, string) error {
	return f.err
}

// TestGetTranscriptRequiresScopeForRedactedView 验证该路径的预期行为。
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

// TestGetTranscriptAllowsRedactedViewWhenAuthorized 验证该路径的预期行为。
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

// TestExportTranscriptRequiresExportPermission 验证该路径的预期行为。
func TestExportTranscriptRequiresExportPermission(t *testing.T) {
	handler := NewHandler(fakeRuntimeService{}, fakeAuthorizer{err: authz.ErrForbidden})
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/executor-sessions/session_1/transcript/export", nil)
	req.Header.Set("X-Actor-Id", "demo-user")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected status %d, got %d", http.StatusForbidden, rec.Code)
	}
}
