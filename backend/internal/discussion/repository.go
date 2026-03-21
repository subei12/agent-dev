package discussion

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateSession(ctx context.Context, cmd CreateSessionCmd) (Session, error) {
	row, err := sqlc.New(r.pool).CreateDiscussionSession(ctx, sqlc.CreateDiscussionSessionParams{
		ID:                 uuid.NewString(),
		MissionID:          cmd.MissionID,
		Topic:              cmd.Topic,
		Status:             "open",
		InitiatedByAgentID: cmd.InitiatedByAgentID,
	})
	if err != nil {
		return Session{}, err
	}
	return sessionFromRow(row), nil
}

func (r *Repository) ListSessions(ctx context.Context, missionID string) ([]Session, error) {
	rows, err := sqlc.New(r.pool).ListDiscussionSessionsByMission(ctx, missionID)
	if err != nil {
		return nil, err
	}
	sessions := make([]Session, 0, len(rows))
	for _, row := range rows {
		sessions = append(sessions, sessionFromRow(row))
	}
	return sessions, nil
}

func (r *Repository) CreateRound(ctx context.Context, cmd CreateRoundCmd) (Round, error) {
	roundNo, err := sqlc.New(r.pool).GetNextDiscussionRoundNumber(ctx, cmd.SessionID)
	if err != nil {
		return Round{}, err
	}

	row, err := sqlc.New(r.pool).CreateDiscussionRound(ctx, sqlc.CreateDiscussionRoundParams{
		ID:                           uuid.NewString(),
		SessionID:                    cmd.SessionID,
		RoundNo:                      roundNo,
		PromptSummary:                cmd.PromptSummary,
		ParticipantAgentIdsJson:      jsonValue(cmd.ParticipantAgentIDs, []byte("[]")),
		Summary:                      cmd.Summary,
		OpenQuestionsJson:            jsonValue(cmd.OpenQuestions, []byte("[]")),
		ConflictsJson:                jsonValue(cmd.Conflicts, []byte("{}")),
		ConclusionJson:               jsonValue(cmd.Conclusion, []byte("{}")),
		AdoptedDocumentVersionIdsJson: jsonValue(cmd.AdoptedDocumentVersionIDs, []byte("[]")),
	})
	if err != nil {
		return Round{}, err
	}
	return roundFromRow(row), nil
}

func sessionFromRow(row sqlc.DiscussionSession) Session {
	return Session{
		ID:                 row.ID,
		MissionID:          row.MissionID,
		Topic:              row.Topic,
		Status:             row.Status,
		InitiatedByAgentID: row.InitiatedByAgentID,
		CreatedAt:          row.CreatedAt.Time,
	}
}

func roundFromRow(row sqlc.DiscussionRound) Round {
	return Round{
		ID:                        row.ID,
		SessionID:                 row.SessionID,
		RoundNo:                   row.RoundNo,
		PromptSummary:             row.PromptSummary,
		ParticipantAgentIDs:       row.ParticipantAgentIdsJson,
		Summary:                   row.Summary,
		OpenQuestions:             row.OpenQuestionsJson,
		Conflicts:                 row.ConflictsJson,
		Conclusion:                row.ConclusionJson,
		AdoptedDocumentVersionIDs: row.AdoptedDocumentVersionIdsJson,
		CreatedAt:                 row.CreatedAt.Time,
	}
}

func jsonValue(value, fallback []byte) []byte {
	if len(value) == 0 {
		return fallback
	}
	return value
}

var _ pgtype.Text
