package archive

import (
	"context"
	"testing"
)

type fakeStore struct{}

func (f fakeStore) BuildArchiveData(context.Context, string) (ArchiveData, error) {
	return ArchiveData{
		AdoptedDocumentVersionIDs: []string{"docv_1"},
		IncludedDecisionIDs:       []string{"decision_1"},
		RuntimeEventSummaryKeys:   []string{"runtime/summary.json"},
	}, nil
}

func (f fakeStore) CreateArchiveRecord(_ context.Context, missionID string, manifest ArchiveManifest) (MissionArchive, error) {
	return MissionArchive{
		ID:        "archive_1",
		MissionID: missionID,
		Status:    "uploaded",
		Hash:      manifest.ManifestHash,
	}, nil
}

func (f fakeStore) GetArchiveByMission(context.Context, string) (MissionArchive, error) {
	return MissionArchive{ID: "archive_1", MissionID: "mission_1", Status: "uploaded"}, nil
}

func TestBuildArchiveManifestIncludesRuntimeSummary(t *testing.T) {
	svc := NewService(fakeStore{})
	manifest, err := svc.BuildManifest(context.Background(), "mission_1")
	if err != nil {
		t.Fatalf("build manifest: %v", err)
	}
	if len(manifest.ObjectKeys) == 0 {
		t.Fatal("expected manifest object keys")
	}
	if len(manifest.IncludedDecisionIDs) == 0 {
		t.Fatal("expected included decision ids")
	}
}
