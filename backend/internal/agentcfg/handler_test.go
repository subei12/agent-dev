package agentcfg

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	items []AgentConfig
}

func (f fakeService) List(context.Context, string) ([]AgentConfig, error) {
	return f.items, nil
}

func (f fakeService) Update(context.Context, UpdateAgentConfigCmd) (AgentConfig, error) {
	return AgentConfig{}, nil
}

func TestListAgentConfigs(t *testing.T) {
	handler := NewHandler(fakeService{
		items: []AgentConfig{
			{ID: "agent_1", Name: "后端 Agent", ExecutorType: "codex_cli", Command: "codex"},
		},
	})

	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/agents", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
