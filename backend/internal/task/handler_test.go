package task

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type fakeService struct {
	board TaskBoard
}

func (f fakeService) CreateBoard(context.Context, CreateBoardCmd) (TaskBoard, error) {
	return TaskBoard{}, nil
}

func (f fakeService) GetBoard(context.Context, string) (TaskBoard, error) {
	return f.board, nil
}

func (f fakeService) Claim(context.Context, ClaimTaskCmd) (TaskClaim, error) {
	return TaskClaim{}, nil
}

func (f fakeService) CreateHandoff(context.Context, CreateHandoffCmd) (TaskHandoff, error) {
	return TaskHandoff{}, nil
}

func (f fakeService) RequestCheckpoint(context.Context, RequestCheckpointCmd) (ReviewCheckpoint, error) {
	return ReviewCheckpoint{}, nil
}

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
