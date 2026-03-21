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

// CreateBoard 创建请求的资源或记录。
func (f *fakeStore) CreateBoard(_ context.Context, missionID, title string) (TaskBoard, error) {
	return TaskBoard{ID: "board_1", MissionID: missionID, Title: title}, nil
}

// CreateTaskItem 创建请求的资源或记录。
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

// GetBoard 返回请求的资源或值。
func (f *fakeStore) GetBoard(_ context.Context, missionID string) (TaskBoard, error) {
	return TaskBoard{ID: "board_1", MissionID: missionID, Title: "Delivery Board"}, nil
}

// ListTaskItems 返回当前查询对应的集合结果。
func (f *fakeStore) ListTaskItems(_ context.Context, _ string) ([]TaskItem, error) {
	items := make([]TaskItem, 0, len(f.taskItems))
	for _, item := range f.taskItems {
		items = append(items, item)
	}
	return items, nil
}

// GetTaskItem 返回请求的资源或值。
func (f *fakeStore) GetTaskItem(_ context.Context, taskItemID string) (TaskItem, error) {
	return f.taskItems[taskItemID], nil
}

// UpdateTaskStatus 更新请求的资源状态。
func (f *fakeStore) UpdateTaskStatus(_ context.Context, taskItemID, status string) (TaskItem, error) {
	item := f.taskItems[taskItemID]
	item.Status = status
	f.taskItems[taskItemID] = item
	return item, nil
}

// CreateClaim 创建请求的资源或记录。
func (f *fakeStore) CreateClaim(_ context.Context, cmd ClaimTaskCmd) (TaskClaim, error) {
	return TaskClaim{ID: "claim_1", TaskItemID: cmd.TaskItemID, AgentID: cmd.AgentID, Status: "active"}, nil
}

// CreateHandoff 创建请求的资源或记录。
func (f *fakeStore) CreateHandoff(_ context.Context, cmd CreateHandoffCmd) (TaskHandoff, error) {
	return TaskHandoff{ID: "handoff_1", TaskItemID: cmd.TaskItemID, ToAdminAgent: cmd.ToAdminAgent, Summary: cmd.Summary, Status: "pending"}, nil
}

// CreateCheckpoint 创建请求的资源或记录。
func (f *fakeStore) CreateCheckpoint(_ context.Context, cmd RequestCheckpointCmd) (ReviewCheckpoint, error) {
	return ReviewCheckpoint{ID: "checkpoint_1", MissionID: cmd.MissionID, TaskItemID: cmd.TaskItemID}, nil
}

// TestCreateHandoffRequiresAdminFallbackWhenNoDownstream 验证该路径的预期行为。
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
