package task

import (
	"context"
	"encoding/json"
	"errors"
)

type Store interface {
	CreateBoard(context.Context, string, string) (TaskBoard, error)
	CreateTaskItem(context.Context, string, CreateTaskItemCmd) (TaskItem, error)
	GetBoard(context.Context, string) (TaskBoard, error)
	ListTaskItems(context.Context, string) ([]TaskItem, error)
	GetTaskItem(context.Context, string) (TaskItem, error)
	UpdateTaskStatus(context.Context, string, string) (TaskItem, error)
	CreateClaim(context.Context, ClaimTaskCmd) (TaskClaim, error)
	CreateHandoff(context.Context, CreateHandoffCmd) (TaskHandoff, error)
	CreateCheckpoint(context.Context, RequestCheckpointCmd) (ReviewCheckpoint, error)
}

type Service interface {
	CreateBoard(context.Context, CreateBoardCmd) (TaskBoard, error)
	GetBoard(context.Context, string) (TaskBoard, error)
	Claim(context.Context, ClaimTaskCmd) (TaskClaim, error)
	CreateHandoff(context.Context, CreateHandoffCmd) (TaskHandoff, error)
	RequestCheckpoint(context.Context, RequestCheckpointCmd) (ReviewCheckpoint, error)
}

type service struct {
	store Store
}

// NewService 创建并返回对应的组件。
func NewService(store Store) Service {
	return &service{store: store}
}

// CreateBoard 创建请求的资源或记录。
func (s *service) CreateBoard(ctx context.Context, cmd CreateBoardCmd) (TaskBoard, error) {
	board, err := s.store.CreateBoard(ctx, cmd.MissionID, cmd.Title)
	if err != nil {
		return TaskBoard{}, err
	}

	for _, item := range cmd.Items {
		created, err := s.store.CreateTaskItem(ctx, board.ID, item)
		if err != nil {
			return TaskBoard{}, err
		}
		board.Items = append(board.Items, created)
	}

	return board, nil
}

// GetBoard 返回请求的资源或值。
func (s *service) GetBoard(ctx context.Context, missionID string) (TaskBoard, error) {
	board, err := s.store.GetBoard(ctx, missionID)
	if err != nil {
		return TaskBoard{}, err
	}

	items, err := s.store.ListTaskItems(ctx, board.ID)
	if err != nil {
		return TaskBoard{}, err
	}
	board.Items = items
	return board, nil
}

// Claim 实现当前函数行为。
func (s *service) Claim(ctx context.Context, cmd ClaimTaskCmd) (TaskClaim, error) {
	if _, err := s.store.UpdateTaskStatus(ctx, cmd.TaskItemID, "claimed"); err != nil {
		return TaskClaim{}, err
	}
	return s.store.CreateClaim(ctx, cmd)
}

// CreateHandoff 创建请求的资源或记录。
func (s *service) CreateHandoff(ctx context.Context, cmd CreateHandoffCmd) (TaskHandoff, error) {
	taskItem, err := s.store.GetTaskItem(ctx, cmd.TaskItemID)
	if err != nil {
		return TaskHandoff{}, err
	}

	if !cmd.ToAdminAgent && cmd.ToAgentID == "" && !hasJSONItems(taskItem.DownstreamTaskIDs) {
		return TaskHandoff{}, errors.New("toAdminAgent is required when task has no downstream")
	}

	handoff, err := s.store.CreateHandoff(ctx, cmd)
	if err != nil {
		return TaskHandoff{}, err
	}

	if _, err := s.store.UpdateTaskStatus(ctx, cmd.TaskItemID, "handoff_pending"); err != nil {
		return TaskHandoff{}, err
	}
	return handoff, nil
}

// RequestCheckpoint 实现当前函数行为。
func (s *service) RequestCheckpoint(ctx context.Context, cmd RequestCheckpointCmd) (ReviewCheckpoint, error) {
	return s.store.CreateCheckpoint(ctx, cmd)
}

// hasJSONItems 实现当前函数行为。
func hasJSONItems(raw json.RawMessage) bool {
	if len(raw) == 0 {
		return false
	}
	var items []string
	if err := json.Unmarshal(raw, &items); err != nil {
		return false
	}
	return len(items) > 0
}
