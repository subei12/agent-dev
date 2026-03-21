package runtime

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Store interface {
	CreateSession(context.Context, StartSessionCmd) (ExecutorSession, error)
	UpsertPresence(context.Context, UpsertPresenceCmd) error
	UpsertMissionRuntime(context.Context, UpsertMissionRuntimeCmd) error
	CreateTranscript(context.Context, CreateTranscriptCmd) (Transcript, error)
	CreateEvent(context.Context, AppendEventCmd) error
	GetTranscriptBySession(context.Context, string) (Transcript, error)
	NextTranscriptSeq(context.Context, string) (int32, error)
	CreateTranscriptEntry(context.Context, string, int32, AppendTranscriptEntryCmd) error
	GetSession(context.Context, string) (ExecutorSession, error)
	UpdateSessionStatus(context.Context, string, string) (ExecutorSession, error)
	UpdateTranscriptStatus(context.Context, string, string) (Transcript, error)
	ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error)
	ListSessionEvents(context.Context, string) ([]RuntimeEvent, error)
	ListTranscriptEntries(context.Context, string) ([]TranscriptEntry, error)
	CreateAccessAudit(context.Context, string, string, string) error
	ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository 创建并返回对应的组件。
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CreateSession 创建请求的资源或记录。
func (r *Repository) CreateSession(ctx context.Context, cmd StartSessionCmd) (ExecutorSession, error) {
	row, err := sqlc.New(r.pool).CreateExecutorSession(ctx, sqlc.CreateExecutorSessionParams{
		ID:                uuid.NewString(),
		MissionID:         cmd.MissionID,
		TaskItemID:        textValue(cmd.TaskItemID),
		AgentID:           cmd.AgentID,
		ExecutorProfileID: cmd.ExecutorProfileID,
		Backend:           cmd.Backend,
		Status:            cmd.Status,
	})
	if err != nil {
		return ExecutorSession{}, err
	}
	return sessionFromRow(row), nil
}

// UpsertPresence 实现当前函数行为。
func (r *Repository) UpsertPresence(ctx context.Context, cmd UpsertPresenceCmd) error {
	_, err := sqlc.New(r.pool).UpsertAgentPresence(ctx, sqlc.UpsertAgentPresenceParams{
		AgentID:                  cmd.AgentID,
		Availability:             cmd.Availability,
		CurrentMissionID:         textValue(cmd.CurrentMissionID),
		CurrentTaskItemID:        textValue(cmd.CurrentTaskItemID),
		CurrentExecutorSessionID: textValue(cmd.CurrentExecutorSessionID),
	})
	return err
}

// UpsertMissionRuntime 实现当前函数行为。
func (r *Repository) UpsertMissionRuntime(ctx context.Context, cmd UpsertMissionRuntimeCmd) error {
	_, err := sqlc.New(r.pool).UpsertMissionAgentRuntime(ctx, sqlc.UpsertMissionAgentRuntimeParams{
		ID:                       uuid.NewString(),
		MissionID:                cmd.MissionID,
		AgentID:                  cmd.AgentID,
		Status:                   cmd.Status,
		CurrentTaskItemID:        textValue(cmd.CurrentTaskItemID),
		CurrentRunID:             textValue(cmd.CurrentRunID),
		CurrentNodeRunID:         textValue(cmd.CurrentNodeRunID),
		CurrentExecutorSessionID: textValue(cmd.CurrentExecutorSessionID),
		StatusSummary:            textValue(cmd.StatusSummary),
		StartedAt:                pgtype.Timestamptz{},
	})
	return err
}

// CreateTranscript 创建请求的资源或记录。
func (r *Repository) CreateTranscript(ctx context.Context, cmd CreateTranscriptCmd) (Transcript, error) {
	row, err := sqlc.New(r.pool).CreateExecutorTranscript(ctx, sqlc.CreateExecutorTranscriptParams{
		ID:                   uuid.NewString(),
		ExecutorSessionID:    cmd.ExecutorSessionID,
		MissionID:            cmd.MissionID,
		AgentID:              cmd.AgentID,
		StorageKind:          cmd.StorageKind,
		RedactedObjectKey:    textValue(cmd.RedactedObjectKey),
		RawEncryptedObjectKey: textValue(cmd.RawEncryptedObjectKey),
		Status:               cmd.Status,
	})
	if err != nil {
		return Transcript{}, err
	}
	return transcriptFromRow(row), nil
}

// CreateEvent 创建请求的资源或记录。
func (r *Repository) CreateEvent(ctx context.Context, cmd AppendEventCmd) error {
	_, err := sqlc.New(r.pool).CreateAgentRuntimeEvent(ctx, sqlc.CreateAgentRuntimeEventParams{
		ID:                uuid.NewString(),
		MissionID:         cmd.MissionID,
		AgentID:           cmd.AgentID,
		ExecutorSessionID: cmd.ExecutorSessionID,
		RunID:             textValue(cmd.RunID),
		NodeRunID:         textValue(cmd.NodeRunID),
		TaskItemID:        textValue(cmd.TaskItemID),
		Level:             cmd.Level,
		Category:          cmd.Category,
		Type:              cmd.Type,
		Title:             cmd.Title,
		Summary:           textValue(cmd.Summary),
		PayloadJson:       jsonValue(cmd.Payload, []byte("{}")),
	})
	return err
}

// GetTranscriptBySession 返回请求的资源或值。
func (r *Repository) GetTranscriptBySession(ctx context.Context, sessionID string) (Transcript, error) {
	row, err := sqlc.New(r.pool).GetExecutorTranscriptBySession(ctx, sessionID)
	if err != nil {
		return Transcript{}, err
	}
	return transcriptFromRow(row), nil
}

// NextTranscriptSeq 实现当前函数行为。
func (r *Repository) NextTranscriptSeq(ctx context.Context, transcriptID string) (int32, error) {
	return sqlc.New(r.pool).GetNextTranscriptEntrySeq(ctx, transcriptID)
}

// CreateTranscriptEntry 创建请求的资源或记录。
func (r *Repository) CreateTranscriptEntry(ctx context.Context, transcriptID string, seq int32, cmd AppendTranscriptEntryCmd) error {
	_, err := sqlc.New(r.pool).CreateExecutorTranscriptEntry(ctx, sqlc.CreateExecutorTranscriptEntryParams{
		ID:           uuid.NewString(),
		TranscriptID: transcriptID,
		Seq:          seq,
		Role:         cmd.Role,
		EntryType:    cmd.EntryType,
		RedactedText: textValue(Redact(cmd.ContentText)),
		ContentJson:  jsonValue(cmd.ContentJSON, []byte("{}")),
		RawObjectKey: textValue(cmd.RawObjectKey),
	})
	return err
}

// GetSession 返回请求的资源或值。
func (r *Repository) GetSession(ctx context.Context, sessionID string) (ExecutorSession, error) {
	row, err := sqlc.New(r.pool).GetExecutorSession(ctx, sessionID)
	if err != nil {
		return ExecutorSession{}, err
	}
	return sessionFromRow(row), nil
}

// UpdateSessionStatus 更新请求的资源状态。
func (r *Repository) UpdateSessionStatus(ctx context.Context, sessionID, status string) (ExecutorSession, error) {
	row, err := sqlc.New(r.pool).UpdateExecutorSessionStatus(ctx, sqlc.UpdateExecutorSessionStatusParams{
		ID:     sessionID,
		Status: status,
	})
	if err != nil {
		return ExecutorSession{}, err
	}
	return sessionFromRow(row), nil
}

// UpdateTranscriptStatus 更新请求的资源状态。
func (r *Repository) UpdateTranscriptStatus(ctx context.Context, transcriptID, status string) (Transcript, error) {
	row, err := sqlc.New(r.pool).UpdateExecutorTranscriptStatus(ctx, sqlc.UpdateExecutorTranscriptStatusParams{
		ID:     transcriptID,
		Status: status,
	})
	if err != nil {
		return Transcript{}, err
	}
	return transcriptFromRow(row), nil
}

// ListMissionRuntimes 返回当前查询对应的集合结果。
func (r *Repository) ListMissionRuntimes(ctx context.Context, missionID string) ([]MissionAgentRuntime, error) {
	rows, err := sqlc.New(r.pool).ListMissionAgentRuntimes(ctx, missionID)
	if err != nil {
		return nil, err
	}
	items := make([]MissionAgentRuntime, 0, len(rows))
	for _, row := range rows {
		items = append(items, missionRuntimeFromRow(row))
	}
	return items, nil
}

// ListSessionEvents 返回当前查询对应的集合结果。
func (r *Repository) ListSessionEvents(ctx context.Context, sessionID string) ([]RuntimeEvent, error) {
	rows, err := sqlc.New(r.pool).ListAgentRuntimeEventsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	items := make([]RuntimeEvent, 0, len(rows))
	for _, row := range rows {
		items = append(items, eventFromRow(row))
	}
	return items, nil
}

// ListTranscriptEntries 返回当前查询对应的集合结果。
func (r *Repository) ListTranscriptEntries(ctx context.Context, sessionID string) ([]TranscriptEntry, error) {
	rows, err := sqlc.New(r.pool).ListExecutorTranscriptEntriesBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	items := make([]TranscriptEntry, 0, len(rows))
	for _, row := range rows {
		items = append(items, transcriptEntryFromRow(row))
	}
	return items, nil
}

// CreateAccessAudit 创建请求的资源或记录。
func (r *Repository) CreateAccessAudit(ctx context.Context, sessionID, actorUserID, accessMode string) error {
	transcript, err := r.GetTranscriptBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	_, err = sqlc.New(r.pool).CreateTranscriptAccessAudit(ctx, sqlc.CreateTranscriptAccessAuditParams{
		ID:          uuid.NewString(),
		TranscriptID: transcript.ID,
		ActorUserID: actorUserID,
		AccessMode:  accessMode,
		Reason:      textValue(""),
	})
	return err
}

// ListAccessAudits 返回当前查询对应的集合结果。
func (r *Repository) ListAccessAudits(ctx context.Context, sessionID string) ([]TranscriptAccessAudit, error) {
	rows, err := sqlc.New(r.pool).ListTranscriptAccessAuditsBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	items := make([]TranscriptAccessAudit, 0, len(rows))
	for _, row := range rows {
		items = append(items, accessAuditFromRow(row))
	}
	return items, nil
}

// sessionFromRow 实现当前函数行为。
func sessionFromRow(row sqlc.ExecutorSession) ExecutorSession {
	return ExecutorSession{
		ID:                row.ID,
		MissionID:         row.MissionID,
		TaskItemID:        stringValue(row.TaskItemID),
		AgentID:           row.AgentID,
		ExecutorProfileID: row.ExecutorProfileID,
		Backend:           row.Backend,
		Status:            row.Status,
		StartedAt:         row.StartedAt.Time,
		EndedAt:           timeValue(row.EndedAt),
	}
}

// transcriptFromRow 实现当前函数行为。
func transcriptFromRow(row sqlc.ExecutorTranscript) Transcript {
	return Transcript{
		ID:                    row.ID,
		ExecutorSessionID:     row.ExecutorSessionID,
		MissionID:             row.MissionID,
		AgentID:               row.AgentID,
		StorageKind:           row.StorageKind,
		RedactedObjectKey:     stringValue(row.RedactedObjectKey),
		RawEncryptedObjectKey: stringValue(row.RawEncryptedObjectKey),
		Status:                row.Status,
		StartedAt:             row.StartedAt.Time,
		EndedAt:               timeValue(row.EndedAt),
	}
}

// missionRuntimeFromRow 实现当前函数行为。
func missionRuntimeFromRow(row sqlc.MissionAgentRuntime) MissionAgentRuntime {
	return MissionAgentRuntime{
		ID:                       row.ID,
		MissionID:                row.MissionID,
		AgentID:                  row.AgentID,
		Status:                   row.Status,
		CurrentTaskItemID:        stringValue(row.CurrentTaskItemID),
		CurrentRunID:             stringValue(row.CurrentRunID),
		CurrentNodeRunID:         stringValue(row.CurrentNodeRunID),
		CurrentExecutorSessionID: stringValue(row.CurrentExecutorSessionID),
		StatusSummary:            stringValue(row.StatusSummary),
		StartedAt:                timeValue(row.StartedAt),
		UpdatedAt:                row.UpdatedAt.Time,
	}
}

// eventFromRow 实现当前函数行为。
func eventFromRow(row sqlc.AgentRuntimeEvent) RuntimeEvent {
	return RuntimeEvent{
		ID:                row.ID,
		MissionID:         row.MissionID,
		AgentID:           row.AgentID,
		ExecutorSessionID: row.ExecutorSessionID,
		RunID:             stringValue(row.RunID),
		NodeRunID:         stringValue(row.NodeRunID),
		TaskItemID:        stringValue(row.TaskItemID),
		Level:             row.Level,
		Category:          row.Category,
		Type:              row.Type,
		Title:             row.Title,
		Summary:           stringValue(row.Summary),
		Payload:           json.RawMessage(row.PayloadJson),
		OccurredAt:        row.OccurredAt.Time,
	}
}

// transcriptEntryFromRow 实现当前函数行为。
func transcriptEntryFromRow(row sqlc.ExecutorTranscriptEntry) TranscriptEntry {
	return TranscriptEntry{
		ID:           row.ID,
		TranscriptID: row.TranscriptID,
		Seq:          row.Seq,
		Role:         row.Role,
		EntryType:    row.EntryType,
		RedactedText: stringValue(row.RedactedText),
		ContentJSON:  json.RawMessage(row.ContentJson),
		RawObjectKey: stringValue(row.RawObjectKey),
		CreatedAt:    row.CreatedAt.Time,
	}
}

// accessAuditFromRow 实现当前函数行为。
func accessAuditFromRow(row sqlc.TranscriptAccessAudit) TranscriptAccessAudit {
	return TranscriptAccessAudit{
		ID:           row.ID,
		TranscriptID: row.TranscriptID,
		ActorUserID:  row.ActorUserID,
		AccessMode:   row.AccessMode,
		Reason:       stringValue(row.Reason),
		CreatedAt:    row.CreatedAt.Time,
	}
}

// jsonValue 实现当前函数行为。
func jsonValue(value json.RawMessage, fallback []byte) []byte {
	if len(value) == 0 {
		return fallback
	}
	return []byte(value)
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

// timeValue 实现当前函数行为。
func timeValue(value pgtype.Timestamptz) time.Time {
	if !value.Valid {
		return time.Time{}
	}
	return value.Time
}
