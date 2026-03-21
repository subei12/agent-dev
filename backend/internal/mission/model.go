package mission

import "time"

type Mission struct {
	ID           string    `json:"id"`
	ProjectID    string    `json:"projectId"`
	Title        string    `json:"title"`
	Description  string    `json:"description,omitempty"`
	SourceType   string    `json:"sourceType"`
	Status       string    `json:"status"`
	AdminAgentID string    `json:"adminAgentId,omitempty"`
	TeamID       string    `json:"teamId,omitempty"`
	RepoBindingID string   `json:"repoBindingId,omitempty"`
	CreatedBy    string    `json:"createdBy"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type MissionDecision struct {
	ID               string    `json:"id"`
	MissionID        string    `json:"missionId"`
	DecidedByAgentID string    `json:"decidedByAgentId"`
	Decision         string    `json:"decision"`
	Summary          string    `json:"summary"`
	CreatedAt        time.Time `json:"createdAt"`
}

type CreateMissionCmd struct {
	ProjectID     string
	Title         string
	Description   string
	SourceType    string
	TeamID        string
	RepoBindingID string
	AdminAgentID  string
	CreatedBy     string
}

type DecideMissionCmd struct {
	MissionID        string
	DecidedByAgentID string
	Decision         string
	Summary          string
}
