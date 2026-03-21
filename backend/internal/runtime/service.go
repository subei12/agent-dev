package runtime

import (
	"context"
	"encoding/json"
)

type Service interface {
	StartSession(context.Context, StartSessionCmd) (ExecutorSession, error)
	AppendEvent(context.Context, AppendEventCmd) error
	AppendTranscriptEntry(context.Context, AppendTranscriptEntryCmd) error
	SealSession(context.Context, string) error
	ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error)
	GetSession(context.Context, string) (ExecutorSession, error)
	ListSessionEvents(context.Context, string) ([]RuntimeEvent, error)
	GetTranscriptView(context.Context, string, string, string) (TranscriptView, error)
	ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error)
}

type service struct {
	store     Store
	publisher EventPublisher
	writer    ObjectWriter
}

type EventPublisher interface {
	PublishToChannel(channel string, data []byte)
}

type ObjectWriter interface {
	PutJSON(context.Context, string, any) error
}

// NewService 创建并返回对应的组件。
func NewService(store Store, publisher ...EventPublisher) Service {
	var selected EventPublisher
	if len(publisher) > 0 {
		selected = publisher[0]
	}
	return &service{
		store:     store,
		publisher: selected,
	}
}

// NewServiceWithDeps 创建并返回对应的组件。
func NewServiceWithDeps(store Store, publisher EventPublisher, writer ObjectWriter) Service {
	return &service{
		store:     store,
		publisher: publisher,
		writer:    writer,
	}
}

// StartSession 创建运行时会话、同步 mission presence，并准备 transcript 存储。
func (s *service) StartSession(ctx context.Context, cmd StartSessionCmd) (ExecutorSession, error) {
	// 1. 先创建不可变的执行会话记录。
	if cmd.Status == "" {
		cmd.Status = "starting"
	}

	session, err := s.store.CreateSession(ctx, cmd)
	if err != nil {
		return ExecutorSession{}, err
	}

	if err := s.store.UpsertPresence(ctx, UpsertPresenceCmd{
		AgentID:                  cmd.AgentID,
		Availability:             "busy",
		CurrentMissionID:         cmd.MissionID,
		CurrentTaskItemID:        cmd.TaskItemID,
		CurrentExecutorSessionID: session.ID,
	}); err != nil {
		return ExecutorSession{}, err
	}

	// 2. 把活动会话同步到 mission runtime 状态，方便前端展示当前执行归属。
	if err := s.store.UpsertMissionRuntime(ctx, UpsertMissionRuntimeCmd{
		MissionID:                cmd.MissionID,
		AgentID:                  cmd.AgentID,
		Status:                   "working",
		CurrentTaskItemID:        cmd.TaskItemID,
		CurrentExecutorSessionID: session.ID,
		StatusSummary:            "session started",
	}); err != nil {
		return ExecutorSession{}, err
	}

	// 3. 创建 transcript 外壳记录，并向订阅方发布初始 session 事件。
	if _, err := s.store.CreateTranscript(ctx, CreateTranscriptCmd{
		ExecutorSessionID: session.ID,
		MissionID:         cmd.MissionID,
		AgentID:           cmd.AgentID,
		StorageKind:       "db_text",
		RedactedObjectKey: transcriptRedactedObjectKey(session.ID),
		Status:            "capturing",
	}); err != nil {
		return ExecutorSession{}, err
	}

	s.publish(cmd.MissionID, session.ID, "session_started")

	return session, nil
}

// AppendEvent 向当前流追加新的事件或转录条目。
func (s *service) AppendEvent(ctx context.Context, cmd AppendEventCmd) error {
	if err := s.store.CreateEvent(ctx, cmd); err != nil {
		return err
	}
	s.publish(cmd.MissionID, cmd.ExecutorSessionID, cmd.Type)
	return nil
}

// AppendTranscriptEntry 向当前流追加新的事件或转录条目。
func (s *service) AppendTranscriptEntry(ctx context.Context, cmd AppendTranscriptEntryCmd) error {
	transcript, err := s.store.GetTranscriptBySession(ctx, cmd.ExecutorSessionID)
	if err != nil {
		return err
	}
	seq, err := s.store.NextTranscriptSeq(ctx, transcript.ID)
	if err != nil {
		return err
	}
	if err := s.store.CreateTranscriptEntry(ctx, transcript.ID, seq, cmd); err != nil {
		return err
	}
	s.publish(transcript.MissionID, cmd.ExecutorSessionID, "transcript_appended")
	return nil
}

// SealSession 完成当前运行时会话的收尾，并持久化脱敏 transcript 产物。
func (s *service) SealSession(ctx context.Context, sessionID string) error {
	// 1. 先在数据库里把 session 和 transcript 标记为 sealed。
	session, err := s.store.UpdateSessionStatus(ctx, sessionID, "completed")
	if err != nil {
		return err
	}
	transcript, err := s.store.GetTranscriptBySession(ctx, sessionID)
	if err != nil {
		return err
	}
	if _, err := s.store.UpdateTranscriptStatus(ctx, transcript.ID, "sealed"); err != nil {
		return err
	}

	// 2. 物化脱敏 transcript 对象，供审计和归档流程使用。
	if err := s.persistRedactedTranscript(ctx, transcript, sessionID); err != nil {
		return err
	}

	// 3. 重置实时 presence 和 mission runtime 状态为已完成快照。
	if err := s.store.UpsertPresence(ctx, UpsertPresenceCmd{
		AgentID:                  session.AgentID,
		Availability:             "online",
		CurrentMissionID:         "",
		CurrentTaskItemID:        "",
		CurrentExecutorSessionID: "",
	}); err != nil {
		return err
	}
	if err := s.store.UpsertMissionRuntime(ctx, UpsertMissionRuntimeCmd{
		MissionID:                session.MissionID,
		AgentID:                  session.AgentID,
		Status:                   "completed",
		CurrentTaskItemID:        session.TaskItemID,
		CurrentExecutorSessionID: session.ID,
		StatusSummary:            "session completed",
	}); err != nil {
		return err
	}

	s.publish(session.MissionID, session.ID, "session_completed")
	return nil
}

// ListMissionRuntimes 返回当前查询对应的集合结果。
func (s *service) ListMissionRuntimes(ctx context.Context, missionID string) ([]MissionAgentRuntime, error) {
	return s.store.ListMissionRuntimes(ctx, missionID)
}

// GetSession 返回请求的资源或值。
func (s *service) GetSession(ctx context.Context, sessionID string) (ExecutorSession, error) {
	return s.store.GetSession(ctx, sessionID)
}

// ListSessionEvents 返回当前查询对应的集合结果。
func (s *service) ListSessionEvents(ctx context.Context, sessionID string) ([]RuntimeEvent, error) {
	return s.store.ListSessionEvents(ctx, sessionID)
}

// GetTranscriptView returns the requested transcript view and records access audits for redacted reads.
func (s *service) GetTranscriptView(ctx context.Context, sessionID, actorUserID, view string) (TranscriptView, error) {
	transcript, err := s.store.GetTranscriptBySession(ctx, sessionID)
	if err != nil {
		return TranscriptView{}, err
	}

	result := TranscriptView{Transcript: transcript}
	if view == "redacted" {
		entries, err := s.store.ListTranscriptEntries(ctx, sessionID)
		if err != nil {
			return TranscriptView{}, err
		}
		result.Entries = entries
		if err := s.store.CreateAccessAudit(ctx, sessionID, actorUserID, "redacted_transcript"); err != nil {
			return TranscriptView{}, err
		}
	}
	return result, nil
}

// ListAccessAudits 返回当前查询对应的集合结果。
func (s *service) ListAccessAudits(ctx context.Context, sessionID string) ([]TranscriptAccessAudit, error) {
	return s.store.ListAccessAudits(ctx, sessionID)
}

// publish 将运行时通知分发到 mission 级和 session 级 SSE 通道。
func (s *service) publish(missionID, sessionID, eventType string) {
	if s.publisher == nil {
		return
	}

	payload, err := json.Marshal(map[string]string{
		"type":      eventType,
		"missionId": missionID,
		"sessionId": sessionID,
	})
	if err != nil {
		return
	}

	s.publisher.PublishToChannel("mission:"+missionID, payload)
	if sessionID != "" {
		s.publisher.PublishToChannel("session:"+sessionID, payload)
	}
}

// persistRedactedTranscript 将 sealed 后的脱敏 transcript 写入对象存储。
func (s *service) persistRedactedTranscript(ctx context.Context, transcript Transcript, sessionID string) error {
	if s.writer == nil || transcript.RedactedObjectKey == "" {
		return nil
	}

	entries, err := s.store.ListTranscriptEntries(ctx, sessionID)
	if err != nil {
		return err
	}

	payload := map[string]any{
		"transcript": transcript,
		"entries":    entries,
	}
	return s.writer.PutJSON(ctx, transcript.RedactedObjectKey, payload)
}

// transcriptRedactedObjectKey 实现当前函数行为。
func transcriptRedactedObjectKey(sessionID string) string {
	return "transcripts/" + sessionID + "/redacted.json"
}
