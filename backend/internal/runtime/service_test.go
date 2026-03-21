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

func (f *fakePublisher) PublishToChannel(channel string, _ []byte) {
	f.channels = append(f.channels, channel)
}

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

func (f *fakeStore) UpsertPresence(context.Context, UpsertPresenceCmd) error { return nil }
func (f *fakeStore) UpsertMissionRuntime(context.Context, UpsertMissionRuntimeCmd) error { return nil }

func (f *fakeStore) CreateTranscript(_ context.Context, cmd CreateTranscriptCmd) (Transcript, error) {
	transcript := Transcript{
		ID:                "transcript_1",
		ExecutorSessionID: cmd.ExecutorSessionID,
		MissionID:         cmd.MissionID,
		AgentID:           cmd.AgentID,
		StorageKind:       cmd.StorageKind,
		Status:            cmd.Status,
	}
	f.transcripts[cmd.ExecutorSessionID] = transcript
	return transcript, nil
}

func (f *fakeStore) CreateEvent(context.Context, AppendEventCmd) error { return nil }

func (f *fakeStore) GetTranscriptBySession(_ context.Context, sessionID string) (Transcript, error) {
	return f.transcripts[sessionID], nil
}

func (f *fakeStore) NextTranscriptSeq(_ context.Context, _ string) (int32, error) {
	return int32(len(f.entries) + 1), nil
}

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

func (f *fakeStore) GetSession(_ context.Context, sessionID string) (ExecutorSession, error) {
	return f.sessions[sessionID], nil
}

func (f *fakeStore) UpdateSessionStatus(_ context.Context, sessionID, status string) (ExecutorSession, error) {
	session := f.sessions[sessionID]
	session.Status = status
	f.sessions[sessionID] = session
	return session, nil
}

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

func (f *fakeStore) ListMissionRuntimes(context.Context, string) ([]MissionAgentRuntime, error) {
	return nil, nil
}

func (f *fakeStore) ListSessionEvents(context.Context, string) ([]RuntimeEvent, error) {
	return nil, nil
}

func (f *fakeStore) ListTranscriptEntries(context.Context, string) ([]TranscriptEntry, error) {
	return f.entries, nil
}

func (f *fakeStore) CreateAccessAudit(context.Context, string, string, string) error { return nil }
func (f *fakeStore) ListAccessAudits(context.Context, string) ([]TranscriptAccessAudit, error) {
	return nil, nil
}

func TestRedactorMasksKnownSecrets(t *testing.T) {
	got := Redact("token=abc123 bearer sk-live-xyz")
	if strings.Contains(got, "abc123") || strings.Contains(got, "sk-live-xyz") {
		t.Fatalf("expected secrets to be redacted, got %q", got)
	}
	if !strings.Contains(got, "[REDACTED]") {
		t.Fatalf("expected redaction marker, got %q", got)
	}
}

func TestSessionLifecycle(t *testing.T) {
	store := &fakeStore{
		sessions:    map[string]ExecutorSession{},
		transcripts: map[string]Transcript{},
	}
	publisher := &fakePublisher{}

	svc := NewService(store, publisher)
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
}
