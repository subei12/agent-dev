package archive

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type ArchiveManifest struct {
	MissionID                string   `json:"missionId"`
	ArchivedAt               string   `json:"archivedAt"`
	AdoptedDocumentVersionIDs []string `json:"adoptedDocumentVersionIds"`
	IncludedApprovalIDs      []string `json:"includedApprovalIds"`
	IncludedDecisionIDs      []string `json:"includedDecisionIds"`
	ObjectKeys               []string `json:"objectKeys"`
	ManifestHash             string   `json:"manifestHash"`
}

type ArchiveData struct {
	AdoptedDocumentVersionIDs []string
	IncludedApprovalIDs       []string
	IncludedDecisionIDs       []string
	RuntimeEventSummaryKeys   []string
}

func BuildManifest(missionID string, data ArchiveData) ArchiveManifest {
	now := time.Now().UTC().Format(time.RFC3339)
	objectKeys := append([]string{}, data.RuntimeEventSummaryKeys...)
	hashInput := missionID + now
	for _, key := range objectKeys {
		hashInput += key
	}
	sum := sha256.Sum256([]byte(hashInput))

	return ArchiveManifest{
		MissionID:                 missionID,
		ArchivedAt:                now,
		AdoptedDocumentVersionIDs: data.AdoptedDocumentVersionIDs,
		IncludedApprovalIDs:       data.IncludedApprovalIDs,
		IncludedDecisionIDs:       data.IncludedDecisionIDs,
		ObjectKeys:                objectKeys,
		ManifestHash:              hex.EncodeToString(sum[:]),
	}
}
