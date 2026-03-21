package discussion

import "context"

type Service interface {
	CreateSession(context.Context, CreateSessionCmd) (Session, error)
	ListSessions(context.Context, string) ([]Session, error)
	CreateRound(context.Context, CreateRoundCmd) (Round, error)
}

type service struct {
	repo *Repository
}

// NewService 创建并返回对应的组件。
func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

// CreateSession 创建请求的资源或记录。
func (s *service) CreateSession(ctx context.Context, cmd CreateSessionCmd) (Session, error) {
	return s.repo.CreateSession(ctx, cmd)
}

// ListSessions 返回当前查询对应的集合结果。
func (s *service) ListSessions(ctx context.Context, missionID string) ([]Session, error) {
	return s.repo.ListSessions(ctx, missionID)
}

// CreateRound 创建请求的资源或记录。
func (s *service) CreateRound(ctx context.Context, cmd CreateRoundCmd) (Round, error) {
	return s.repo.CreateRound(ctx, cmd)
}
