package document

import "time"

type Document struct {
	ID                      string    `json:"id"`
	MissionID               string    `json:"missionId"`
	Kind                    string    `json:"kind"`
	Title                   string    `json:"title"`
	CurrentAdoptedVersionID string    `json:"currentAdoptedVersionId,omitempty"`
	CreatedAt               time.Time `json:"createdAt"`
	UpdatedAt               time.Time `json:"updatedAt"`
}

type DocumentVersion struct {
	ID                string    `json:"id"`
	DocumentID        string    `json:"documentId"`
	Version           int32     `json:"version"`
	ParentVersionID   string    `json:"parentVersionId,omitempty"`
	Status            string    `json:"status"`
	ContentFormat     string    `json:"contentFormat"`
	StorageKind       string    `json:"storageKind"`
	ObjectKey         string    `json:"objectKey,omitempty"`
	ContentHash       string    `json:"contentHash"`
	ContentText       string    `json:"contentText"`
	ProducedByAgentID string    `json:"producedByAgentId,omitempty"`
	ProducedByRunID   string    `json:"producedByRunId,omitempty"`
	SourceRoundID     string    `json:"sourceRoundId,omitempty"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CreateDocumentCmd struct {
	MissionID string
	Kind      string
	Title     string
}

type CreateVersionCmd struct {
	DocumentID        string
	ParentVersionID   string
	ContentFormat     string
	StorageKind       string
	ObjectKey         string
	ContentHash       string
	ContentText       string
	ProducedByAgentID string
	ProducedByRunID   string
	SourceRoundID     string
}
