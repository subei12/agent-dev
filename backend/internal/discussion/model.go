package discussion

import "time"

type Session struct {
	ID                 string    `json:"id"`
	MissionID          string    `json:"missionId"`
	Topic              string    `json:"topic"`
	Status             string    `json:"status"`
	InitiatedByAgentID string    `json:"initiatedByAgentId"`
	CreatedAt          time.Time `json:"createdAt"`
}

type Round struct {
	ID                        string    `json:"id"`
	SessionID                 string    `json:"sessionId"`
	RoundNo                   int32     `json:"roundNo"`
	PromptSummary             string    `json:"promptSummary"`
	ParticipantAgentIDs       []byte    `json:"participantAgentIdsJson"`
	Summary                   string    `json:"summary"`
	OpenQuestions             []byte    `json:"openQuestionsJson"`
	Conflicts                 []byte    `json:"conflictsJson"`
	Conclusion                []byte    `json:"conclusionJson"`
	AdoptedDocumentVersionIDs []byte    `json:"adoptedDocumentVersionIdsJson"`
	CreatedAt                 time.Time `json:"createdAt"`
}

type CreateSessionCmd struct {
	MissionID          string
	Topic              string
	InitiatedByAgentID string
}

type CreateRoundCmd struct {
	SessionID                 string
	PromptSummary             string
	ParticipantAgentIDs       []byte
	Summary                   string
	OpenQuestions             []byte
	Conflicts                 []byte
	Conclusion                []byte
	AdoptedDocumentVersionIDs []byte
}
