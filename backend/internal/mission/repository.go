package mission

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Repository struct {
	queries *sqlc.Queries
}

// NewRepository 创建并返回对应的组件。
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: sqlc.New(pool),
	}
}

// CreateMission 创建请求的资源或记录。
func (r *Repository) CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error) {
	row, err := r.queries.CreateMission(ctx, sqlc.CreateMissionParams{
		ID:            uuid.NewString(),
		ProjectID:     cmd.ProjectID,
		Title:         cmd.Title,
		Description:   textValue(cmd.Description),
		SourceType:    cmd.SourceType,
		Status:        "discussion",
		AdminAgentID:  cmd.AdminAgentID,
		TeamID:        textValue(cmd.TeamID),
		RepoBindingID: textValue(cmd.RepoBindingID),
		CreatedBy:     cmd.CreatedBy,
	})
	if err != nil {
		return Mission{}, err
	}
	return missionFromRow(row), nil
}

// GetMission 返回请求的资源或值。
func (r *Repository) GetMission(ctx context.Context, projectID, missionID string) (Mission, error) {
	row, err := r.queries.GetMission(ctx, sqlc.GetMissionParams{
		ID:        missionID,
		ProjectID: projectID,
	})
	if err != nil {
		return Mission{}, err
	}
	return missionFromRow(row), nil
}

// CreateDecision 创建请求的资源或记录。
func (r *Repository) CreateDecision(ctx context.Context, cmd DecideMissionCmd) (MissionDecision, error) {
	row, err := r.queries.CreateMissionDecision(ctx, sqlc.CreateMissionDecisionParams{
		ID:                          uuid.NewString(),
		MissionID:                   cmd.MissionID,
		DecidedByAgentID:            cmd.DecidedByAgentID,
		Decision:                    cmd.Decision,
		Summary:                     cmd.Summary,
		RelatedTaskItemID:           textValue(""),
		RelatedDocumentVersionIdsJson: []byte("[]"),
		RelatedRepoCandidateIdsJson:   []byte("[]"),
	})
	if err != nil {
		return MissionDecision{}, err
	}
	return decisionFromRow(row), nil
}

// UpdateMissionStatus 更新请求的资源状态。
func (r *Repository) UpdateMissionStatus(ctx context.Context, missionID, status string) (Mission, error) {
	row, err := r.queries.UpdateMissionStatus(ctx, sqlc.UpdateMissionStatusParams{
		ID:     missionID,
		Status: status,
	})
	if err != nil {
		return Mission{}, err
	}
	return missionFromRow(row), nil
}

// missionFromRow 实现当前函数行为。
func missionFromRow(row sqlc.Mission) Mission {
	return Mission{
		ID:            row.ID,
		ProjectID:     row.ProjectID,
		Title:         row.Title,
		Description:   stringValue(row.Description),
		SourceType:    row.SourceType,
		Status:        row.Status,
		AdminAgentID:  row.AdminAgentID,
		TeamID:        stringValue(row.TeamID),
		RepoBindingID: stringValue(row.RepoBindingID),
		CreatedBy:     row.CreatedBy,
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}

// decisionFromRow 实现当前函数行为。
func decisionFromRow(row sqlc.MissionDecision) MissionDecision {
	return MissionDecision{
		ID:               row.ID,
		MissionID:        row.MissionID,
		DecidedByAgentID: row.DecidedByAgentID,
		Decision:         row.Decision,
		Summary:          row.Summary,
		CreatedAt:        row.CreatedAt.Time,
	}
}

// textValue 实现当前函数行为。
func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

// stringValue 实现当前函数行为。
func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
