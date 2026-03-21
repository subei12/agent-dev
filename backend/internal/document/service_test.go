package document

import (
	"context"
	"testing"
)

type fakeStore struct {
	documents map[string]Document
	versions  map[string]DocumentVersion
}

// CreateDocument 创建请求的资源或记录。
func (f *fakeStore) CreateDocument(_ context.Context, cmd CreateDocumentCmd) (Document, error) {
	document := Document{
		ID:        "doc_" + cmd.Title,
		MissionID: cmd.MissionID,
		Kind:      cmd.Kind,
		Title:     cmd.Title,
	}
	f.documents[document.ID] = document
	return document, nil
}

// GetDocument 返回请求的资源或值。
func (f *fakeStore) GetDocument(_ context.Context, missionID, documentID string) (Document, error) {
	document := f.documents[documentID]
	if document.MissionID != missionID {
		document.MissionID = missionID
	}
	return document, nil
}

// ListDocuments 返回当前查询对应的集合结果。
func (f *fakeStore) ListDocuments(_ context.Context, missionID string) ([]Document, error) {
	documents := make([]Document, 0, len(f.documents))
	for _, document := range f.documents {
		if document.MissionID == missionID {
			documents = append(documents, document)
		}
	}
	return documents, nil
}

// NextVersionNumber 实现当前函数行为。
func (f *fakeStore) NextVersionNumber(_ context.Context, documentID string) (int32, error) {
	var maxVersion int32
	for _, version := range f.versions {
		if version.DocumentID == documentID && version.Version > maxVersion {
			maxVersion = version.Version
		}
	}
	return maxVersion + 1, nil
}

// CreateVersion 创建请求的资源或记录。
func (f *fakeStore) CreateVersion(_ context.Context, cmd CreateVersionCmd, versionNumber int32) (DocumentVersion, error) {
	version := DocumentVersion{
		ID:            cmd.ContentHash,
		DocumentID:    cmd.DocumentID,
		Version:       versionNumber,
		Status:        "proposed",
		ContentFormat: cmd.ContentFormat,
		StorageKind:   cmd.StorageKind,
		ContentHash:   cmd.ContentHash,
		ContentText:   cmd.ContentText,
	}
	f.versions[version.ID] = version
	return version, nil
}

// GetVersion 返回请求的资源或值。
func (f *fakeStore) GetVersion(_ context.Context, versionID string) (DocumentVersion, error) {
	return f.versions[versionID], nil
}

// ListVersions 返回当前查询对应的集合结果。
func (f *fakeStore) ListVersions(_ context.Context, documentID string) ([]DocumentVersion, error) {
	versions := make([]DocumentVersion, 0, len(f.versions))
	for _, version := range f.versions {
		if version.DocumentID == documentID {
			versions = append(versions, version)
		}
	}
	return versions, nil
}

// UpdateVersionStatus 更新请求的资源状态。
func (f *fakeStore) UpdateVersionStatus(_ context.Context, versionID, status string) (DocumentVersion, error) {
	version := f.versions[versionID]
	version.Status = status
	f.versions[versionID] = version
	return version, nil
}

// ClearAdoptedVersions 实现当前函数行为。
func (f *fakeStore) ClearAdoptedVersions(_ context.Context, documentID string) error {
	for id, version := range f.versions {
		if version.DocumentID == documentID && version.Status == "adopted" {
			version.Status = "proposed"
			f.versions[id] = version
		}
	}
	return nil
}

// SetCurrentAdoptedVersion 实现当前函数行为。
func (f *fakeStore) SetCurrentAdoptedVersion(ctx context.Context, documentID, versionID string) (Document, error) {
	_ = f.ClearAdoptedVersions(ctx, documentID)
	version := f.versions[versionID]
	version.Status = "adopted"
	f.versions[versionID] = version

	document := f.documents[documentID]
	document.CurrentAdoptedVersionID = versionID
	f.documents[documentID] = document
	return document, nil
}

// TestAdoptDocumentVersionSetsSingleAdoptedVersion 验证该路径的预期行为。
func TestAdoptDocumentVersionSetsSingleAdoptedVersion(t *testing.T) {
	store := &fakeStore{
		documents: map[string]Document{},
		versions:  map[string]DocumentVersion{},
	}

	svc := NewService(store)
	doc, err := svc.CreateDocument(context.Background(), CreateDocumentCmd{
		MissionID: "mission_1",
		Kind:      "architecture",
		Title:     "Architecture",
	})
	if err != nil {
		t.Fatalf("create document: %v", err)
	}

	v1, err := svc.CreateVersion(context.Background(), CreateVersionCmd{
		DocumentID:   doc.ID,
		ContentText:  "v1",
		ContentHash:  "hash-v1",
		ContentFormat:"md",
		StorageKind:  "db_text",
	})
	if err != nil {
		t.Fatalf("create v1: %v", err)
	}

	v2, err := svc.CreateVersion(context.Background(), CreateVersionCmd{
		DocumentID:   doc.ID,
		ContentText:  "v2",
		ContentHash:  "hash-v2",
		ContentFormat:"md",
		StorageKind:  "db_text",
	})
	if err != nil {
		t.Fatalf("create v2: %v", err)
	}

	if _, err := svc.AdoptVersion(context.Background(), doc.ID, v1.ID); err != nil {
		t.Fatalf("adopt v1: %v", err)
	}

	adopted, err := svc.AdoptVersion(context.Background(), doc.ID, v2.ID)
	if err != nil {
		t.Fatalf("adopt v2: %v", err)
	}

	if adopted.ID != v2.ID {
		t.Fatalf("expected adopted version %q, got %q", v2.ID, adopted.ID)
	}

	gotDoc, err := svc.GetDocument(context.Background(), doc.MissionID, doc.ID)
	if err != nil {
		t.Fatalf("get document: %v", err)
	}
	if gotDoc.CurrentAdoptedVersionID != v2.ID {
		t.Fatalf("expected current adopted version %q, got %q", v2.ID, gotDoc.CurrentAdoptedVersionID)
	}

	versions, err := svc.ListVersions(context.Background(), doc.ID)
	if err != nil {
		t.Fatalf("list versions: %v", err)
	}

	var adoptedCount int
	for _, version := range versions {
		if version.Status == "adopted" {
			adoptedCount++
			if version.ID != v2.ID {
				t.Fatalf("expected only %q to be adopted, got %q", v2.ID, version.ID)
			}
		}
	}

	if adoptedCount != 1 {
		t.Fatalf("expected exactly 1 adopted version, got %d", adoptedCount)
	}
}
