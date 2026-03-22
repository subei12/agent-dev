package mission

import "context"

type Service interface {
	CreateMission(context.Context, CreateMissionCmd) (Mission, error)
	ListMissions(context.Context, string) ([]Mission, error)
	GetMission(context.Context, string, string) (Mission, error)
	Decide(context.Context, DecideMissionCmd) (MissionDecision, error)
}

type service struct {
	repo *Repository
}

// NewService 创建并返回对应的组件。
func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

// CreateMission 创建请求的资源或记录。
func (s *service) CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error) {
	if cmd.SourceType == "" {
		cmd.SourceType = "project_task"
	}
	if cmd.CreatedBy == "" {
		cmd.CreatedBy = "system"
	}
	return s.repo.CreateMission(ctx, cmd)
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
