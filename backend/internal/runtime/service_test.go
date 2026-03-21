package runtime

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
)

type fakeStore struct {
	sessions map[string]ExecutorSession
	transcripts map[string]Transcript
	entries []TranscriptEntry
}

type fakePublisher struct {
	channels []string
}

type fakeWriter struct {
	objectKeys []string
}

// PublishToChannel 将当前事件或负载发布到目标位置。
func (f *fakePublisher) PublishToChannel(channel string, _ []byte) {
	f.channels = append(f.channels, channel)
}

// PutJSON 实现当前函数行为。
func (f *fakeWriter) PutJSON(_ context.Context, objectKey string, _ any) error {
	f.objectKeys = append(f.objectKeys, objectKey)
	return nil
}

// CreateSession 创建请求的资源或记录。
func (f *fakeStore) CreateSession(_ context.Context, cmd StartSessionCmd) (ExecutorSession, error) {
	session := ExecutorSession{
		ID:                "session_1",
		MissionID:         cmd.MissionID,
		TaskItemID:        cmd.TaskItemID,
		AgentID:           cmd.AgentID,
		ExecutorProfileID: cmd.ExecutorProfileID,
		Backend:           cmd.Backend,
		Status:            cmd.Status,
	}
	f.sessions[session.ID] = session
	return session, nil
}

// UpsertPresence 实现当前函数行为。
func (f *fakeStore) UpsertPresence(context.Context, UpsertPresenceCmd) error { return nil }
// UpsertMissionRuntime 实现当前函数行为。
func (f *fakeStore) UpsertMissionRuntime(context.Context, UpsertMissionRuntimeCmd) error { return nil }

// CreateTranscript 创建请求的资源或记录。
func (f *fakeStore) CreateTranscript(_ context.Context, cmd CreateTranscriptCmd) (Transcript, error) {
	transcript := Transcript{
		ID:                "transcript_1",
		ExecutorSessionID: cmd.ExecutorSessionID,
		MissionID:         cmd.MissionID,
		AgentID:           cmd.AgentID,
		StorageKind:       cmd.StorageKind,
		RedactedObjectKey: cmd.RedactedObjectKey,
		Status:            cmd.Status,
	}
	f.transcripts[cmd.ExecutorSessionID] = transcript
	return transcript, nil
}

// CreateEvent 创建请求的资源或记录。
func (f *fakeStore) CreateEvent(context.Context, AppendEventCmd) error { return nil }

// GetTranscriptBySession 返回请求的资源或值。
func (f *fakeStore) GetTranscriptBySession(_ context.Context, sessionID string) (Transcript, error) {
	return f.transcripts[sessionID], nil
}

// NextTranscriptSeq 实现当前函数行为。
func (f *fakeStore) NextTranscriptSeq(_ context.Context, _ string) (int32, error) {
	return int32(len(f.entries) + 1), nil
}

// CreateTranscriptEntry 创建请求的资源或记录。
func (f *fakeStore) CreateTranscriptEntry(_ context.Context, transcriptID string, seq int32, cmd AppendTranscriptEntryCmd) error {
	f.entries = append(f.entries, TranscriptEntry{
		ID:           "entry_1",
		TranscriptID: transcriptID,
		Seq:          seq,
		Role:         cmd.Role,
		EntryType:    cmd.EntryType,
		RedactedText: Redact(cmd.ContentText),
		ContentJSON:  json.RawMessage("{}"),
	})
	return nil
}

// GetSession 返回请求的资源或值。
func (f *fakeStore) GetSession(_ context.Context, sessionID string) (ExecutorSession, error) {
	return f.sessions[sessionID], nil
}

// UpdateSessionStatus 更新请求的资源状态。
func (f *fakeStore) UpdateSessionStatus(_ context.Context, sessionID, status string) (ExecutorSession, error) {
	session := f.sessions[sessionID]
	session.Status = status
	f.sessions[sessionID] = session
	return session, nil
}

// UpdateTranscriptStatus 更新请求的资源状态。
func (f *fakeStore) UpdateTranscriptStatus(_ context.Context, transcriptID, status string) (Transcript, error) {
	for key, transcript := range f.transcripts {
		if transcript.ID == transcriptID {
			transcript.Status = status
			f.transcripts[key] = transcript
			return transcript, nil
		}
	}
	return Transcript{}, nil
}

// ListMissionRuntimes 返回当前查询对应的集合结果。
func (f *fakeStore) ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error) {
	return nil, nil
}

// ListSessionEvents 返回当前查询对应的集合结果。
func (f *fakeStore) ListSessionEvents(context.Context, string) ([]RuntimeEvent, error) {
	return nil, nil
}

// ListTranscriptEntries 返回当前查询对应的集合结果。
func (f *fakeStore) ListTranscriptEntries(context.Context, string) ([]TranscriptEntry, error) {
	return f.entries, nil
}

// CreateAccessAudit 创建请求的资源或记录。
func (f *fakeStore) CreateAccessAudit(context.Context, string, string, string) error { return nil }
// ListAccessAudits 返回当前查询对应的集合结果。
func (f *fakeStore) ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error) {
	return nil, nil
}

// TestRedactorMasksKnownSecrets 验证该路径的预期行为。
func TestRedactorMasksKnownSecrets(t *testing.T) {
	got := Redact("token=abc123 bearer sk-live-xyz")
	if strings.Contains(got, "abc123") || strings.Contains(got, "sk-live-xyz") {
		t.Fatalf("expected secrets to be redacted, got %q", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %q", got)
	}
}

// TestSessionLifecycle 验证该路径的预期行为。
func TestSessionLifecycle(t *testing.T) {
	store := &fakeStore{
		sessions:    map[string]ExecutorSession{},
		transcripts: map[string]Transcript{},
	}
	publisher := &fakePublisher{}
	writer := &fakeWriter{}

	svc := NewServiceWithDeps(store, publisher, writer)
	session, err := svc.StartSession(context.Background(), StartSessionCmd{
		MissionID:          "mission_1",
		AgentID:            "agent_1",
		ExecutorProfileID:  "exec_1",
		Backend:            "codex_cli",
		Status:             "starting",
	})
	if err != nil {
		t.Fatalf("start session: %v", err)
	}

	if err := svc.AppendEvent(context.Background(), AppendEventCmd{
		MissionID:         "mission_1",
		AgentID:           "agent_1",
		ExecutorSessionID: session.ID,
		Level:             "info",
		Category:          "session",
		Type:              "session_started",
		Title:             "Session started",
	}); err != nil {
		t.Fatalf("append event: %v", err)
	}

	if err := svc.AppendTranscriptEntry(context.Background(), AppendTranscriptEntryCmd{
		ExecutorSessionID: session.ID,
		Role:              "assistant",
		EntryType:         "message",
		ContentText:       "token=abc123",
	}); err != nil {
		t.Fatalf("append transcript: %v", err)
	}

	if err := svc.SealSession(context.Background(), session.ID); err != nil {
		t.Fatalf("seal session: %v", err)
	}

	if len(publisher.channels) == 0 {
		t.Fatal("expected runtime events to publish SSE notifications")
	}
	if len(writer.objectKeys) != 1 {
		t.Fatalf("expected 1 transcript object write, got %d", len(writer.objectKeys))
	}
}
