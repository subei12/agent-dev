package task

import (
	"encoding/json"
	"time"
)

type TaskBoard struct {
	ID        string     `json:"id"`
	MissionID string     `json:"missionId"`
	Title     string     `json:"title"`
	Items     []TaskItem `json:"items,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

type TaskItem struct {
	ID                      string          `json:"id"`
	BoardID                 string          `json:"boardId"`
	Title                   string          `json:"title"`
	Type                    string          `json:"type"`
	Status                  string          `json:"status"`
	AssignedAgentID         string          `json:"assignedAgentId,omitempty"`
	UpstreamTaskIDs         json.RawMessage `json:"upstreamTaskIds"`
	DownstreamTaskIDs       json.RawMessage `json:"downstreamTaskIds"`
	InputDocumentVersionIDs json.RawMessage `json:"inputDocumentVersionIds"`
	InputRepoCandidateIDs   json.RawMessage `json:"inputRepoCandidateIds"`
	DefinitionOfDone        json.RawMessage `json:"definitionOfDone"`
	CreatedAt               time.Time       `json:"createdAt"`
	UpdatedAt               time.Time       `json:"updatedAt"`
}

type TaskClaim struct {
	ID          string    `json:"id"`
	TaskItemID  string    `json:"taskItemId"`
	AgentID     string    `json:"agentId"`
	Status      string    `json:"status"`
	ClaimReason string    `json:"claimReason,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

type TaskHandoff struct {
	ID                       string          `json:"id"`
	TaskItemID               string          `json:"taskItemId"`
	FromAgentID              string          `json:"fromAgentId"`
	ToAgentID                string          `json:"toAgentId,omitempty"`
	ToAdminAgent             bool            `json:"toAdminAgent"`
	Summary                  string          `json:"summary"`
	OutputDocumentVersionIDs json.RawMessage `json:"outputDocumentVersionIds"`
	OutputRepoCandidateIDs   json.RawMessage `json:"outputRepoCandidateIds"`
	OutputArtifactVersionIDs json.RawMessage `json:"outputArtifactVersionIds"`
	ValidationSummary        string          `json:"validationSummary,omitempty"`
	RiskSummary              string          `json:"riskSummary,omitempty"`
	RecommendedNextAction    string          `json:"recommendedNextAction,omitempty"`
	Status                   string          `json:"status"`
	CreatedAt                time.Time       `json:"createdAt"`
}

type ReviewCheckpoint struct {
	ID                       string          `json:"id"`
	MissionID                string          `json:"missionId"`
	TaskItemID               string          `json:"taskItemId,omitempty"`
	Kind                     string          `json:"kind"`
	RequestedByAgentID       string          `json:"requestedByAgentId"`
	AssignedAgentID          string          `json:"assignedAgentId,omitempty"`
	Status                   string          `json:"status"`
	Summary                  string          `json:"summary,omitempty"`
	LinkedDocumentVersionIDs json.RawMessage `json:"linkedDocumentVersionIds"`
	LinkedRepoCandidateIDs   json.RawMessage `json:"linkedRepoCandidateIds"`
	CreatedAt                time.Time       `json:"createdAt"`
}

type CreateBoardCmd struct {
	MissionID string
	Title     string
	Items     []CreateTaskItemCmd
}

type CreateTaskItemCmd struct {
	Title                   string          `json:"title"`
	Type                    string          `json:"type"`
	AssignedAgentID         string          `json:"assignedAgentId,omitempty"`
	UpstreamTaskIDs         json.RawMessage `json:"upstreamTaskIds"`
	DownstreamTaskIDs       json.RawMessage `json:"downstreamTaskIds"`
	InputDocumentVersionIDs json.RawMessage `json:"inputDocumentVersionIds"`
	InputRepoCandidateIDs   json.RawMessage `json:"inputRepoCandidateIds"`
	DefinitionOfDone        json.RawMessage `json:"definitionOfDone"`
}

type ClaimTaskCmd struct {
	TaskItemID   string
	AgentID      string
	ClaimReason  string
}

type CreateHandoffCmd struct {
	TaskItemID               string
	FromAgentID              string
	ToAgentID                string
	ToAdminAgent             bool
	Summary                  string
	OutputDocumentVersionIDs json.RawMessage
	OutputRepoCandidateIDs   json.RawMessage
	OutputArtifactVersionIDs json.RawMessage
	ValidationSummary        string
	RiskSummary              string
	RecommendedNextAction    string
}

type RequestCheckpointCmd struct {
	MissionID                string
	TaskItemID               string
	Kind                     string
	RequestedByAgentID       string
	AssignedAgentID          string
	Summary                  string
	LinkedDocumentVersionIDs json.RawMessage
	LinkedRepoCandidateIDs   json.RawMessage
}
