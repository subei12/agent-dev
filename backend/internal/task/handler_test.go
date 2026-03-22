package task

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	board TaskBoard
}

// CreateBoard 创建请求的资源或记录。
func (f fakeService) CreateBoard(context.Context, CreateBoardCmd) (TaskBoard, error) {
	return TaskBoard{}, nil
}

// AddTask 创建请求的资源或记录。
func (f fakeService) AddTask(context.Context, string, CreateTaskItemCmd) (TaskItem, error) {
	return TaskItem{ID: "task_new", Title: "新增任务", Type: "code", Status: "todo"}, nil
}

// GetBoard 返回请求的资源或值。
func (f fakeService) GetBoard(context.Context, string) (TaskBoard, error) {
	return f.board, nil
}

// Claim 实现当前函数行为。
func (f fakeService) Claim(context.Context, ClaimTaskCmd) (TaskClaim, error) {
	return TaskClaim{}, nil
}

// CreateHandoff 创建请求的资源或记录。
func (f fakeService) CreateHandoff(context.Context, CreateHandoffCmd) (TaskHandoff, error) {
	return TaskHandoff{}, nil
}

// RequestCheckpoint 实现当前函数行为。
func (f fakeService) RequestCheckpoint(context.Context, RequestCheckpointCmd) (ReviewCheckpoint, error) {
	return ReviewCheckpoint{}, nil
}

// TestGetTaskBoard 验证该路径的预期行为。
func TestGetTaskBoard(t *testing.T) {
	handler := NewHandler(fakeService{
		board: TaskBoard{
			ID:        "board_1",
			MissionID: "mission_1",
			Title:     "Delivery Board",
			Items: []TaskItem{
				{ID: "task_1", Title: "Implement runtime observability", Type: "code", Status: "in_progress"},
			},
		},
	})

	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/projects/proj_1/missions/mission_1/task-board", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}

// TestCreateTask 验证新增任务接口能正常返回。
func TestCreateTask(t *testing.T) {
	handler := NewHandler(fakeService{
		board: TaskBoard{
			ID:        "board_1",
			MissionID: "mission_1",
			Title:     "Delivery Board",
		},
	})
	router := chi.NewRouter()
	handler.RegisterRoutes(router)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/projects/proj_1/missions/mission_1/tasks",
		strings.NewReader(`{"title":"新增任务","type":"code"}`),
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}
