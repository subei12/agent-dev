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

// CreateMission 创建请求的资源或记录。
func (f fakeService) CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error) {
	return f.createMission(ctx, cmd)
}

// ListMissions 返回当前查询对应的集合结果。
func (f fakeService) ListMissions(context.Context, string) ([]Mission, error) {
	return []Mission{
		{ID: "mission_1", Title: "演示任务", Status: "implementation"},
	}, nil
}

// GetMission 返回请求的资源或值。
func (f fakeService) GetMission(context.Context, string, string) (Mission, error) {
	return Mission{}, nil
}

// Decide 实现当前函数行为。
func (f fakeService) Decide(context.Context, DecideMissionCmd) (MissionDecision, error) {
	return MissionDecision{}, nil
}

// TestCreateMission 验证该路径的预期行为。
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

// TestListMissions 验证 Mission 列表接口能正常返回。
func TestListMissions(t *testing.T) {
	handler := NewHandler(fakeService{
		createMission: func(_ context.Context, cmd CreateMissionCmd) (Mission, error) {
			return Mission{ID: "mission_1", ProjectID: cmd.ProjectID, Title: cmd.Title, Status: "discussion"}, nil
		},
	})
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/missions", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
