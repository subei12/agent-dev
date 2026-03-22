package mission

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/your-org/agent-platform/internal/discussion"
	"github.com/your-org/agent-platform/internal/document"
	"github.com/your-org/agent-platform/internal/task"
)

type fakeMissionRepo struct {
	createdMission Mission
}

// CreateMission 创建请求的资源或记录。
func (f *fakeMissionRepo) CreateMission(_ context.Context, cmd CreateMissionCmd) (Mission, error) {
	f.createdMission = Mission{
		ID:            "mission_1",
		ProjectID:     cmd.ProjectID,
		Title:         cmd.Title,
		Description:   cmd.Description,
		SourceType:    cmd.SourceType,
		Status:        "discussion",
		AdminAgentID:  cmd.AdminAgentID,
		TeamID:        cmd.TeamID,
		RepoBindingID: cmd.RepoBindingID,
		CreatedBy:     cmd.CreatedBy,
	}
	return f.createdMission, nil
}

// ListMissions 返回当前查询对应的集合结果。
func (f *fakeMissionRepo) ListMissions(context.Context, string) ([]Mission, error) {
	return nil, nil
}

// GetMission 返回请求的资源或值。
func (f *fakeMissionRepo) GetMission(context.Context, string, string) (Mission, error) {
	return Mission{}, nil
}

// CreateDecision 创建请求的资源或记录。
func (f *fakeMissionRepo) CreateDecision(context.Context, DecideMissionCmd) (MissionDecision, error) {
	return MissionDecision{}, nil
}

// UpdateMissionStatus 更新请求的资源状态。
func (f *fakeMissionRepo) UpdateMissionStatus(context.Context, string, string) (Mission, error) {
	return Mission{}, nil
}

type fakeBootstrapper struct {
	mission Mission
	err     error
}

// Bootstrap 初始化 Mission 的自动协作骨架。
func (f *fakeBootstrapper) Bootstrap(_ context.Context, mission Mission) error {
	f.mission = mission
	return f.err
}

type fakeDiscussionService struct {
	sessions []discussion.CreateSessionCmd
}

// CreateSession 创建请求的资源或记录。
func (f *fakeDiscussionService) CreateSession(_ context.Context, cmd discussion.CreateSessionCmd) (discussion.Session, error) {
	f.sessions = append(f.sessions, cmd)
	return discussion.Session{
		ID:                 "discussion_1",
		MissionID:          cmd.MissionID,
		Topic:              cmd.Topic,
		Status:             "open",
		InitiatedByAgentID: cmd.InitiatedByAgentID,
	}, nil
}

type fakeDocumentService struct {
	createDocumentCmd document.CreateDocumentCmd
	createVersionCmd  document.CreateVersionCmd
	adoptDocumentID   string
	adoptVersionID    string
}

// CreateDocument 创建请求的资源或记录。
func (f *fakeDocumentService) CreateDocument(_ context.Context, cmd document.CreateDocumentCmd) (document.Document, error) {
	f.createDocumentCmd = cmd
	return document.Document{
		ID:        "doc_1",
		MissionID: cmd.MissionID,
		Kind:      cmd.Kind,
		Title:     cmd.Title,
	}, nil
}

// CreateVersion 创建请求的资源或记录。
func (f *fakeDocumentService) CreateVersion(_ context.Context, cmd document.CreateVersionCmd) (document.DocumentVersion, error) {
	f.createVersionCmd = cmd
	return document.DocumentVersion{
		ID:            "docv_1",
		DocumentID:    cmd.DocumentID,
		Version:       1,
		Status:        "proposed",
		ContentFormat: cmd.ContentFormat,
		StorageKind:   cmd.StorageKind,
		ContentHash:   cmd.ContentHash,
		ContentText:   cmd.ContentText,
	}, nil
}

// AdoptVersion 实现当前函数行为。
func (f *fakeDocumentService) AdoptVersion(_ context.Context, documentID, versionID string) (document.DocumentVersion, error) {
	f.adoptDocumentID = documentID
	f.adoptVersionID = versionID
	return document.DocumentVersion{
		ID:         versionID,
		DocumentID: documentID,
		Status:     "adopted",
	}, nil
}

type fakeTaskService struct {
	boardCmd    task.CreateBoardCmd
	claimCmd    task.ClaimTaskCmd
}

// CreateBoard 创建请求的资源或记录。
func (f *fakeTaskService) CreateBoard(_ context.Context, cmd task.CreateBoardCmd) (task.TaskBoard, error) {
	f.boardCmd = cmd
	items := make([]task.TaskItem, 0, len(cmd.Items))
	for index, item := range cmd.Items {
		items = append(items, task.TaskItem{
			ID:              "task_" + item.Title,
			Title:           item.Title,
			Type:            item.Type,
			Status:          "todo",
			AssignedAgentID: item.AssignedAgentID,
			BoardID:         "board_1",
		})
		if index == 0 {
			items[index].ID = "task_plan"
		}
	}
	return task.TaskBoard{
		ID:        "board_1",
		MissionID: cmd.MissionID,
		Title:     cmd.Title,
		Items:     items,
	}, nil
}

// Claim 创建请求的资源或记录。
func (f *fakeTaskService) Claim(_ context.Context, cmd task.ClaimTaskCmd) (task.TaskClaim, error) {
	f.claimCmd = cmd
	return task.TaskClaim{
		ID:         "claim_1",
		TaskItemID: cmd.TaskItemID,
		AgentID:    cmd.AgentID,
		Status:     "active",
	}, nil
}

// TestCreateMissionTriggersBootstrap 验证创建 Mission 时会调用自动引导器。
func TestCreateMissionTriggersBootstrap(t *testing.T) {
	repo := &fakeMissionRepo{}
	bootstrapper := &fakeBootstrapper{}

	svc := NewService(repo, bootstrapper)
	mission, err := svc.CreateMission(context.Background(), CreateMissionCmd{
		ProjectID:    "proj_1",
		Title:        "实现自动开发流程",
		Description:  "用户只提交需求，平台自动推进。",
		AdminAgentID: "agent_admin",
	})
	if err != nil {
		t.Fatalf("create mission: %v", err)
	}

	if mission.ID == "" {
		t.Fatal("expected created mission id")
	}
	if bootstrapper.mission.ID != mission.ID {
		t.Fatalf("expected bootstrap mission %q, got %q", mission.ID, bootstrapper.mission.ID)
	}
}

// TestCreateMissionReturnsBootstrapError 验证自动引导失败时会返回错误。
func TestCreateMissionReturnsBootstrapError(t *testing.T) {
	repo := &fakeMissionRepo{}
	bootstrapper := &fakeBootstrapper{err: errors.New("bootstrap failed")}

	svc := NewService(repo, bootstrapper)
	_, err := svc.CreateMission(context.Background(), CreateMissionCmd{
		ProjectID:    "proj_1",
		Title:        "实现自动开发流程",
		AdminAgentID: "agent_admin",
	})
	if err == nil {
		t.Fatal("expected bootstrap error")
	}
}

// TestDefaultMissionBootstrapperCreatesDiscussionDocumentAndTasks 验证默认引导器会创建讨论、方案文档和内部任务。
func TestDefaultMissionBootstrapperCreatesDiscussionDocumentAndTasks(t *testing.T) {
	discussionSvc := &fakeDiscussionService{}
	documentSvc := &fakeDocumentService{}
	taskSvc := &fakeTaskService{}

	bootstrapper := NewDefaultBootstrapper(discussionSvc, documentSvc, taskSvc)
	err := bootstrapper.Bootstrap(context.Background(), Mission{
		ID:           "mission_1",
		ProjectID:    "proj_1",
		Title:        "实现多 Agent 自动开发流程",
		Description:  "平台要自动从设计推进到开发和测试。",
		AdminAgentID: "agent_admin",
	})
	if err != nil {
		t.Fatalf("bootstrap mission: %v", err)
	}

	if len(discussionSvc.sessions) != 1 {
		t.Fatalf("expected 1 discussion session, got %d", len(discussionSvc.sessions))
	}
	if discussionSvc.sessions[0].Topic != "需求澄清与方案收敛" {
		t.Fatalf("unexpected discussion topic: %s", discussionSvc.sessions[0].Topic)
	}

	if documentSvc.createDocumentCmd.Title != "统一实施方案" {
		t.Fatalf("unexpected document title: %s", documentSvc.createDocumentCmd.Title)
	}
	if documentSvc.adoptVersionID != "docv_1" {
		t.Fatalf("expected adopted version docv_1, got %s", documentSvc.adoptVersionID)
	}
	if documentSvc.createVersionCmd.ProducedByAgentID != "agent_admin" {
		t.Fatalf("expected producedByAgentID agent_admin, got %s", documentSvc.createVersionCmd.ProducedByAgentID)
	}

	if taskSvc.boardCmd.Title != "内部执行任务" {
		t.Fatalf("unexpected board title: %s", taskSvc.boardCmd.Title)
	}
	if len(taskSvc.boardCmd.Items) != 4 {
		t.Fatalf("expected 4 bootstrap tasks, got %d", len(taskSvc.boardCmd.Items))
	}

	gotDocIDs := decodeTaskIDs(t, taskSvc.boardCmd.Items[0].InputDocumentVersionIDs)
	if len(gotDocIDs) != 1 || gotDocIDs[0] != "docv_1" {
		t.Fatalf("expected bootstrap task to reference adopted document version, got %v", gotDocIDs)
	}

	if taskSvc.boardCmd.Items[0].AssignedAgentID != "agent_admin" {
		t.Fatalf("expected design task assigned to admin agent, got %s", taskSvc.boardCmd.Items[0].AssignedAgentID)
	}
	if taskSvc.boardCmd.Items[1].Type != "code" {
		t.Fatalf("expected second task type code, got %s", taskSvc.boardCmd.Items[1].Type)
	}
	if taskSvc.boardCmd.Items[2].Type != "test" {
		t.Fatalf("expected third task type test, got %s", taskSvc.boardCmd.Items[2].Type)
	}
	if taskSvc.boardCmd.Items[3].Type != "review" {
		t.Fatalf("expected fourth task type review, got %s", taskSvc.boardCmd.Items[3].Type)
	}
	if taskSvc.claimCmd.TaskItemID != "task_plan" {
		t.Fatalf("expected first claimed task task_plan, got %s", taskSvc.claimCmd.TaskItemID)
	}
	if taskSvc.claimCmd.AgentID != "agent_admin" {
		t.Fatalf("expected auto claim assigned to agent_admin, got %s", taskSvc.claimCmd.AgentID)
	}
}

// decodeTaskIDs 解析任务引用中的 JSON ID 列表。
func decodeTaskIDs(t *testing.T, raw json.RawMessage) []string {
	t.Helper()

	var items []string
	if err := json.Unmarshal(raw, &items); err != nil {
		t.Fatalf("decode ids: %v", err)
	}
	return items
}
