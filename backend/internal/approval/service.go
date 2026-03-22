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
	ListByMission(context.Context, string) ([]Approval, error)
}

type service struct {
	queries *sqlc.Queries
}

// NewService 创建并返回对应的组件。
func NewService(pool *pgxpool.Pool) Service {
	return &service{queries: sqlc.New(pool)}
}

// Create 创建请求的资源或记录。
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

// Get 返回请求的资源或值。
func (s *service) Get(ctx context.Context, approvalID string) (Approval, error) {
	row, err := s.queries.GetApproval(ctx, approvalID)
	if err != nil {
		return Approval{}, err
	}
	return approvalFromRow(row), nil
}

// ListByMission 返回当前 Mission 下的审批列表。
func (s *service) ListByMission(ctx context.Context, missionID string) ([]Approval, error) {
	rows, err := s.queries.ListApprovalsByMission(ctx, textValue(missionID))
	if err != nil {
		return nil, err
	}

	items := make([]Approval, 0, len(rows))
	for _, row := range rows {
		items = append(items, approvalFromRow(row))
	}
	return items, nil
}

// approvalFromRow 实现当前函数行为。
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
