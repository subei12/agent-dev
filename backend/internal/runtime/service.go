package runtime

import "context"

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
	store Store
}

func NewService(store Store) Service {
	return &service{store: store}
}

func (s *service) StartSession(ctx context.Context, cmd StartSessionCmd) (ExecutorSession, error) {
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

	if _, err := s.store.CreateTranscript(ctx, CreateTranscriptCmd{
		ExecutorSessionID: session.ID,
		MissionID:         cmd.MissionID,
		AgentID:           cmd.AgentID,
		StorageKind:       "db_text",
		Status:            "capturing",
	}); err != nil {
		return ExecutorSession{}, err
	}

	return session, nil
}

func (s *service) AppendEvent(ctx context.Context, cmd AppendEventCmd) error {
	return s.store.CreateEvent(ctx, cmd)
}

func (s *service) AppendTranscriptEntry(ctx context.Context, cmd AppendTranscriptEntryCmd) error {
	transcript, err := s.store.GetTranscriptBySession(ctx, cmd.ExecutorSessionID)
	if err != nil {
		return err
	}
	seq, err := s.store.NextTranscriptSeq(ctx, transcript.ID)
	if err != nil {
		return err
	}
	return s.store.CreateTranscriptEntry(ctx, transcript.ID, seq, cmd)
}

func (s *service) SealSession(ctx context.Context, sessionID string) error {
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
	if err := s.store.UpsertPresence(ctx, UpsertPresenceCmd{
		AgentID:                  session.AgentID,
		Availability:             "online",
		CurrentMissionID:         "",
		CurrentTaskItemID:        "",
		CurrentExecutorSessionID: "",
	}); err != nil {
		return err
	}
	return s.store.UpsertMissionRuntime(ctx, UpsertMissionRuntimeCmd{
		MissionID:                session.MissionID,
		AgentID:                  session.AgentID,
		Status:                   "completed",
		CurrentTaskItemID:        session.TaskItemID,
		CurrentExecutorSessionID: session.ID,
		StatusSummary:            "session completed",
	})
}

func (s *service) ListMissionRuntimes(ctx context.Context, missionID string) ([]MissionAgentRuntime, error) {
	return s.store.ListMissionRuntimes(ctx, missionID)
}

func (s *service) GetSession(ctx context.Context, sessionID string) (ExecutorSession, error) {
	return s.store.GetSession(ctx, sessionID)
}

func (s *service) ListSessionEvents(ctx context.Context, sessionID string) ([]RuntimeEvent, error) {
	return s.store.ListSessionEvents(ctx, sessionID)
}

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

func (s *service) ListAccessAudits(ctx context.Context, sessionID string) ([]TranscriptAccessAudit, error) {
	return s.store.ListAccessAudits(ctx, sessionID)
}
