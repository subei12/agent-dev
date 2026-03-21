package approval

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Approval struct {
	ID         string          `json:"id"`
	MissionID  string          `json:"missionId,omitempty"`
	Action     string          `json:"action"`
	SubjectType string         `json:"subjectType"`
	SubjectID  string          `json:"subjectId"`
	Intent     json.RawMessage `json:"intentSnapshot"`
	IntentHash string          `json:"intentHash"`
	Status     string          `json:"status"`
	CreatedBy  string          `json:"createdBy"`
}

type CreateApprovalCmd struct {
	MissionID   string
	Action      string
	SubjectType string
	SubjectID   string
	Intent      json.RawMessage
	CreatedBy   string
}

type Service interface {
	Create(context.Context, CreateApprovalCmd) (Approval, error)
	Get(context.Context, string) (Approval, error)
}

type service struct {
	queries *sqlc.Queries
}

func NewService(pool *pgxpool.Pool) Service {
	return &service{queries: sqlc.New(pool)}
}

func (s *service) Create(ctx context.Context, cmd CreateApprovalCmd) (Approval, error) {
	intent := cmd.Intent
	if len(intent) == 0 {
		intent = json.RawMessage("{}")
	}
	hash := sha256.Sum256(intent)

	row, err := s.queries.CreateApproval(ctx, sqlc.CreateApprovalParams{
		ID:               uuid.NewString(),
		MissionID:        textValue(cmd.MissionID),
		RunID:            pgtype.Text{},
		NodeRunID:        pgtype.Text{},
		Action:           cmd.Action,
		SubjectType:      cmd.SubjectType,
		SubjectID:        cmd.SubjectID,
		IntentSnapshotJson: []byte(intent),
		IntentHash:       hex.EncodeToString(hash[:]),
		Status:           "pending",
		Comment:          pgtype.Text{},
		CreatedBy:        cmd.CreatedBy,
		DecidedBy:        pgtype.Text{},
	})
	if err != nil {
		return Approval{}, err
	}
	return approvalFromRow(row), nil
}

func (s *service) Get(ctx context.Context, approvalID string) (Approval, error) {
	row, err := s.queries.GetApproval(ctx, approvalID)
	if err != nil {
		return Approval{}, err
	}
	return approvalFromRow(row), nil
}

func approvalFromRow(row sqlc.Approval) Approval {
	return Approval{
		ID:          row.ID,
		MissionID:   stringValue(row.MissionID),
		Action:      row.Action,
		SubjectType: row.SubjectType,
		SubjectID:   row.SubjectID,
		Intent:      json.RawMessage(row.IntentSnapshotJson),
		IntentHash:  row.IntentHash,
		Status:      row.Status,
		CreatedBy:   row.CreatedBy,
	}
}

func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}
