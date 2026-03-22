package mission

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/your-org/agent-platform/internal/discussion"
	"github.com/your-org/agent-platform/internal/document"
	"github.com/your-org/agent-platform/internal/task"
)

type Service interface {
	CreateMission(context.Context, CreateMissionCmd) (Mission, error)
	ListMissions(context.Context, string) ([]Mission, error)
	GetMission(context.Context, string, string) (Mission, error)
	Decide(context.Context, DecideMissionCmd) (MissionDecision, error)
}

type repository interface {
	CreateMission(context.Context, CreateMissionCmd) (Mission, error)
	ListMissions(context.Context, string) ([]Mission, error)
	GetMission(context.Context, string, string) (Mission, error)
	CreateDecision(context.Context, DecideMissionCmd) (MissionDecision, error)
	UpdateMissionStatus(context.Context, string, string) (Mission, error)
}

type Bootstrapper interface {
	Bootstrap(context.Context, Mission) error
}

type discussionService interface {
	CreateSession(context.Context, discussion.CreateSessionCmd) (discussion.Session, error)
}

type documentService interface {
	CreateDocument(context.Context, document.CreateDocumentCmd) (document.Document, error)
	CreateVersion(context.Context, document.CreateVersionCmd) (document.DocumentVersion, error)
	AdoptVersion(context.Context, string, string) (document.DocumentVersion, error)
}

type taskService interface {
	CreateBoard(context.Context, task.CreateBoardCmd) (task.TaskBoard, error)
}

type service struct {
	repo        repository
	bootstrapper Bootstrapper
}

// NewService 创建并返回对应的组件。
func NewService(repo repository, bootstrapper ...Bootstrapper) Service {
	var selected Bootstrapper
	if len(bootstrapper) > 0 {
		selected = bootstrapper[0]
	}
	return &service{repo: repo, bootstrapper: selected}
}

// CreateMission 创建请求的资源或记录。
func (s *service) CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error) {
	if cmd.SourceType == "" {
		cmd.SourceType = "project_task"
	}
	if cmd.CreatedBy == "" {
		cmd.CreatedBy = "system"
	}

	mission, err := s.repo.CreateMission(ctx, cmd)
	if err != nil {
		return Mission{}, err
	}
	if s.bootstrapper != nil {
		if err := s.bootstrapper.Bootstrap(ctx, mission); err != nil {
			return Mission{}, err
		}
	}
	return mission, nil
}

// ListMissions 返回指定项目下的 Mission 列表。
func (s *service) ListMissions(ctx context.Context, projectID string) ([]Mission, error) {
	return s.repo.ListMissions(ctx, projectID)
}

// GetMission 返回请求的资源或值。
func (s *service) GetMission(ctx context.Context, projectID, missionID string) (Mission, error) {
	return s.repo.GetMission(ctx, projectID, missionID)
}

// Decide 实现当前函数行为。
func (s *service) Decide(ctx context.Context, cmd DecideMissionCmd) (MissionDecision, error) {
	decision, err := s.repo.CreateDecision(ctx, cmd)
	if err != nil {
		return MissionDecision{}, err
	}

	if nextStatus := statusFromDecision(cmd.Decision); nextStatus != "" {
		if _, err := s.repo.UpdateMissionStatus(ctx, cmd.MissionID, nextStatus); err != nil {
			return MissionDecision{}, err
		}
	}

	return decision, nil
}

// statusFromDecision 实现当前函数行为。
func statusFromDecision(decision string) string {
	switch decision {
	case "approve_plan":
		return "planning"
	case "enter_implementation":
		return "implementation"
	case "complete_mission":
		return "completed"
	default:
		return ""
	}
}

type defaultBootstrapper struct {
	discussions discussionService
	documents   documentService
	tasks       taskService
}

// NewDefaultBootstrapper 创建 Mission 自动协作骨架引导器。
func NewDefaultBootstrapper(discussions discussionService, documents documentService, tasks taskService) Bootstrapper {
	return &defaultBootstrapper{
		discussions: discussions,
		documents:   documents,
		tasks:       tasks,
	}
}

// Bootstrap 为新建 Mission 自动生成讨论、方案文档和内部执行任务。
func (b *defaultBootstrapper) Bootstrap(ctx context.Context, mission Mission) error {
	// 1. 先创建第一轮需求澄清讨论，让管理员 Agent 有统一收敛入口。
	if _, err := b.discussions.CreateSession(ctx, discussion.CreateSessionCmd{
		MissionID:          mission.ID,
		Topic:              "需求澄清与方案收敛",
		InitiatedByAgentID: mission.AdminAgentID,
	}); err != nil {
		return err
	}

	// 2. 基于用户需求生成统一实施方案文档，并直接采纳为当前共识版本。
	planDocument, err := b.documents.CreateDocument(ctx, document.CreateDocumentCmd{
		MissionID: mission.ID,
		Kind:      "implementation_plan",
		Title:     "统一实施方案",
	})
	if err != nil {
		return err
	}

	planVersion, err := b.documents.CreateVersion(ctx, document.CreateVersionCmd{
		DocumentID:        planDocument.ID,
		ContentFormat:     "md",
		StorageKind:       "db_text",
		ContentHash:       fmt.Sprintf("bootstrap-%s-plan", mission.ID),
		ContentText:       buildPlanDraft(mission),
		ProducedByAgentID: mission.AdminAgentID,
	})
	if err != nil {
		return err
	}

	adoptedVersion, err := b.documents.AdoptVersion(ctx, planDocument.ID, planVersion.ID)
	if err != nil {
		return err
	}

	// 3. 按设计 -> 开发 -> 测试 -> 复核的顺序生成内部执行任务，供多 Agent 接力处理。
	_, err = b.tasks.CreateBoard(ctx, task.CreateBoardCmd{
		MissionID: mission.ID,
		Title:     "内部执行任务",
		Items: []task.CreateTaskItemCmd{
			{
				Title:                   "管理员 Agent 收敛方案",
				Type:                    "design",
				AssignedAgentID:         mission.AdminAgentID,
				DownstreamTaskIDs:       mustJSON([]string{"开发实现与自测"}),
				InputDocumentVersionIDs: mustJSON([]string{adoptedVersion.ID}),
				DefinitionOfDone:        mustJSON(map[string]string{"done": "统一方案已采纳"}),
			},
			{
				Title:                   "开发实现与自测",
				Type:                    "code",
				AssignedAgentID:         "agent_backend",
				UpstreamTaskIDs:         mustJSON([]string{"管理员 Agent 收敛方案"}),
				DownstreamTaskIDs:       mustJSON([]string{"测试与验收"}),
				InputDocumentVersionIDs: mustJSON([]string{adoptedVersion.ID}),
				DefinitionOfDone:        mustJSON(map[string]string{"done": "完成实现与基础自测"}),
			},
			{
				Title:                   "测试与验收",
				Type:                    "test",
				AssignedAgentID:         "agent_backend",
				UpstreamTaskIDs:         mustJSON([]string{"开发实现与自测"}),
				DownstreamTaskIDs:       mustJSON([]string{"管理员复核与结项"}),
				InputDocumentVersionIDs: mustJSON([]string{adoptedVersion.ID}),
				DefinitionOfDone:        mustJSON(map[string]string{"done": "测试通过并形成验收结论"}),
			},
			{
				Title:                   "管理员复核与结项",
				Type:                    "review",
				AssignedAgentID:         mission.AdminAgentID,
				UpstreamTaskIDs:         mustJSON([]string{"测试与验收"}),
				InputDocumentVersionIDs: mustJSON([]string{adoptedVersion.ID}),
				DefinitionOfDone:        mustJSON(map[string]string{"done": "管理员确认完成或追加检查"}),
			},
		},
	})
	return err
}

// buildPlanDraft 生成初始实施方案草稿内容。
func buildPlanDraft(mission Mission) string {
	return fmt.Sprintf("# %s\n\n## 原始需求\n\n%s\n\n## 默认执行路径\n\n1. 管理员 Agent 收敛需求与方案。\n2. 开发 Agent 完成实现与自测。\n3. 测试 Agent 进行测试与验收。\n4. 管理员 Agent 决策是否完成、是否追加检查。", mission.Title, mission.Description)
}

// mustJSON 把任务 bootstrap 输入编码为 JSON RawMessage。
func mustJSON(value any) json.RawMessage {
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	return data
}
