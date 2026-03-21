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

func NewService(repo *Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateSession(ctx context.Context, cmd CreateSessionCmd) (Session, error) {
	return s.repo.CreateSession(ctx, cmd)
}

func (s *service) ListSessions(ctx context.Context, missionID string) ([]Session, error) {
	return s.repo.ListSessions(ctx, missionID)
}

func (s *service) CreateRound(ctx context.Context, cmd CreateRoundCmd) (Round, error) {
	return s.repo.CreateRound(ctx, cmd)
}
