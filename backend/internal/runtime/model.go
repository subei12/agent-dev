package runtime

import (
	"encoding/json"
	"time"
)

type AgentPresence struct {
	AgentID                  string    `json:"agentId"`
	Availability             string    `json:"availability"`
	CurrentMissionID         string    `json:"currentMissionId,omitempty"`
	CurrentTaskItemID        string    `json:"currentTaskItemId,omitempty"`
	CurrentExecutorSessionID string    `json:"currentExecutorSessionId,omitempty"`
	LastHeartbeatAt          time.Time `json:"lastHeartbeatAt"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

type MissionAgentRuntime struct {
	ID                       string    `json:"id"`
	MissionID                string    `json:"missionId"`
	AgentID                  string    `json:"agentId"`
	Status                   string    `json:"status"`
	CurrentTaskItemID        string    `json:"currentTaskItemId,omitempty"`
	CurrentRunID             string    `json:"currentRunId,omitempty"`
	CurrentNodeRunID         string    `json:"currentNodeRunId,omitempty"`
	CurrentExecutorSessionID string    `json:"currentExecutorSessionId,omitempty"`
	StatusSummary            string    `json:"statusSummary,omitempty"`
	StartedAt                time.Time `json:"startedAt,omitempty"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

type ExecutorSession struct {
	ID                string    `json:"id"`
	MissionID         string    `json:"missionId"`
	TaskItemID        string    `json:"taskItemId,omitempty"`
	AgentID           string    `json:"agentId"`
	ExecutorProfileID string    `json:"executorProfileId"`
	Backend           string    `json:"backend"`
	Status            string    `json:"status"`
	StartedAt         time.Time `json:"startedAt"`
	EndedAt           time.Time `json:"endedAt,omitempty"`
}

type RuntimeEvent struct {
	ID                string          `json:"id"`
	MissionID         string          `json:"missionId"`
	AgentID           string          `json:"agentId"`
	ExecutorSessionID string          `json:"executorSessionId"`
	RunID             string          `json:"runId,omitempty"`
	NodeRunID         string          `json:"nodeRunId,omitempty"`
	TaskItemID        string          `json:"taskItemId,omitempty"`
	Level             string          `json:"level"`
	Category          string          `json:"category"`
	Type              string          `json:"type"`
	Title             string          `json:"title"`
	Summary           string          `json:"summary,omitempty"`
	Payload           json.RawMessage `json:"payload"`
	OccurredAt        time.Time       `json:"occurredAt"`
}

type Transcript struct {
	ID                  string    `json:"id"`
	ExecutorSessionID   string    `json:"executorSessionId"`
	MissionID           string    `json:"missionId"`
	AgentID             string    `json:"agentId"`
	StorageKind         string    `json:"storageKind"`
	RedactedObjectKey   string    `json:"redactedObjectKey,omitempty"`
	RawEncryptedObjectKey string  `json:"rawEncryptedObjectKey,omitempty"`
	Status              string    `json:"status"`
	StartedAt           time.Time `json:"startedAt"`
	EndedAt             time.Time `json:"endedAt,omitempty"`
}

type TranscriptEntry struct {
	ID           string          `json:"id"`
	TranscriptID string          `json:"transcriptId"`
	Seq          int32           `json:"seq"`
	Role         string          `json:"role"`
	EntryType    string          `json:"entryType"`
	RedactedText string          `json:"redactedText,omitempty"`
	ContentJSON  json.RawMessage `json:"contentJson"`
	RawObjectKey string          `json:"rawObjectKey,omitempty"`
	CreatedAt    time.Time       `json:"createdAt"`
}

type TranscriptAccessAudit struct {
	ID          string    `json:"id"`
	TranscriptID string   `json:"transcriptId"`
	ActorUserID string    `json:"actorUserId"`
	AccessMode  string    `json:"accessMode"`
	Reason      string    `json:"reason,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TranscriptView struct {
	Transcript Transcript        `json:"transcript"`
	Entries    []TranscriptEntry `json:"entries,omitempty"`
}

type StartSessionCmd struct {
	MissionID         string
	TaskItemID        string
	AgentID           string
	ExecutorProfileID string
	Backend           string
	Status            string
}

type AppendEventCmd struct {
	MissionID         string
	AgentID           string
	ExecutorSessionID string
	RunID             string
	NodeRunID         string
	TaskItemID        string
	Level             string
	Category          string
	Type              string
	Title             string
	Summary           string
	Payload           json.RawMessage
}

type AppendTranscriptEntryCmd struct {
	ExecutorSessionID string
	Role              string
	EntryType         string
	ContentText       string
	ContentJSON       json.RawMessage
	RawObjectKey      string
}

type UpsertPresenceCmd struct {
	AgentID                  string
	Availability             string
	CurrentMissionID         string
	CurrentTaskItemID        string
	CurrentExecutorSessionID string
}

type UpsertMissionRuntimeCmd struct {
	MissionID                string
	AgentID                  string
	Status                   string
	CurrentTaskItemID        string
	CurrentRunID             string
	CurrentNodeRunID         string
	CurrentExecutorSessionID string
	StatusSummary            string
}

type CreateTranscriptCmd struct {
	ExecutorSessionID    string
	MissionID            string
	AgentID              string
	StorageKind          string
	RedactedObjectKey    string
	RawEncryptedObjectKey string
	Status               string
}
