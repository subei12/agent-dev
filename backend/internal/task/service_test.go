package task

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type fakeStore struct {
	taskItems map[string]TaskItem
}

func (f *fakeStore) CreateBoard(_ context.Context, missionID, title string) (TaskBoard, error) {
	return TaskBoard{ID: "board_1", MissionID: missionID, Title: title}, nil
}

func (f *fakeStore) CreateTaskItem(_ context.Context, boardID string, cmd CreateTaskItemCmd) (TaskItem, error) {
	item := TaskItem{
		ID:                cmd.Title,
		BoardID:           boardID,
		Title:             cmd.Title,
		Type:              cmd.Type,
		Status:            "todo",
		DownstreamTaskIDs: cmd.DownstreamTaskIDs,
	}
	f.taskItems[item.ID] = item
	return item, nil
}

func (f *fakeStore) GetTaskItem(_ context.Context, taskItemID string) (TaskItem, error) {
	return f.taskItems[taskItemID], nil
}

func (f *fakeStore) UpdateTaskStatus(_ context.Context, taskItemID, status string) (TaskItem, error) {
	item := f.taskItems[taskItemID]
	item.Status = status
	f.taskItems[taskItemID] = item
	return item, nil
}

func (f *fakeStore) CreateClaim(_ context.Context, cmd ClaimTaskCmd) (TaskClaim, error) {
	return TaskClaim{ID: "claim_1", TaskItemID: cmd.TaskItemID, AgentID: cmd.AgentID, Status: "active"}, nil
}

func (f *fakeStore) CreateHandoff(_ context.Context, cmd CreateHandoffCmd) (TaskHandoff, error) {
	return TaskHandoff{ID: "handoff_1", TaskItemID: cmd.TaskItemID, ToAdminAgent: cmd.ToAdminAgent, Summary: cmd.Summary, Status: "pending"}, nil
}

func (f *fakeStore) CreateCheckpoint(_ context.Context, cmd RequestCheckpointCmd) (ReviewCheckpoint, error) {
	return ReviewCheckpoint{ID: "checkpoint_1", MissionID: cmd.MissionID, TaskItemID: cmd.TaskItemID}, nil
}

func TestCreateHandoffRequiresAdminFallbackWhenNoDownstream(t *testing.T) {
	store := &fakeStore{
		taskItems: map[string]TaskItem{
			"task_1": {
				ID:                "task_1",
				DownstreamTaskIDs: json.RawMessage(`[]`),
			},
		},
	}

	svc := NewService(store)
	_, err := svc.CreateHandoff(context.Background(), CreateHandoffCmd{
		TaskItemID: "task_1",
		FromAgentID: "agent_dev",
		Summary: "done",
	})
	if err == nil {
		t.Fatal("expected error when no downstream and no admin fallback")
	}
	if !strings.Contains(err.Error(), "toAdminAgent") {
		t.Fatalf("expected toAdminAgent error, got %v", err)
	}
}
