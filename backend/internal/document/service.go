package document

import (
	"context"
	"errors"
)

type Service interface {
	CreateDocument(context.Context, CreateDocumentCmd) (Document, error)
	GetDocument(context.Context, string, string) (Document, error)
	ListDocuments(context.Context, string) ([]Document, error)
	CreateVersion(context.Context, CreateVersionCmd) (DocumentVersion, error)
	GetVersion(context.Context, string) (DocumentVersion, error)
	ListVersions(context.Context, string) ([]DocumentVersion, error)
	AdoptVersion(context.Context, string, string) (DocumentVersion, error)
}

type service struct {
	store Store
}

// NewService 创建并返回对应的组件。
func NewService(store Store) Service {
	return &service{store: store}
}

// CreateDocument 创建请求的资源或记录。
func (s *service) CreateDocument(ctx context.Context, cmd CreateDocumentCmd) (Document, error) {
	if cmd.Kind == "" {
		cmd.Kind = "other"
	}
	return s.store.CreateDocument(ctx, cmd)
}

// GetDocument 返回请求的资源或值。
func (s *service) GetDocument(ctx context.Context, missionID, documentID string) (Document, error) {
	return s.store.GetDocument(ctx, missionID, documentID)
}

// ListDocuments 返回当前查询对应的集合结果。
func (s *service) ListDocuments(ctx context.Context, missionID string) ([]Document, error) {
	return s.store.ListDocuments(ctx, missionID)
}

// CreateVersion 创建请求的资源或记录。
func (s *service) CreateVersion(ctx context.Context, cmd CreateVersionCmd) (DocumentVersion, error) {
	if cmd.ContentFormat == "" {
		cmd.ContentFormat = "md"
	}
	if cmd.StorageKind == "" {
		cmd.StorageKind = "db_text"
	}

	nextVersion, err := s.store.NextVersionNumber(ctx, cmd.DocumentID)
	if err != nil {
		return DocumentVersion{}, err
	}
	return s.store.CreateVersion(ctx, cmd, nextVersion)
}

// GetVersion 返回请求的资源或值。
func (s *service) GetVersion(ctx context.Context, versionID string) (DocumentVersion, error) {
	return s.store.GetVersion(ctx, versionID)
}

// ListVersions 返回当前查询对应的集合结果。
func (s *service) ListVersions(ctx context.Context, documentID string) ([]DocumentVersion, error) {
	return s.store.ListVersions(ctx, documentID)
}

// AdoptVersion 实现当前函数行为。
func (s *service) AdoptVersion(ctx context.Context, documentID, versionID string) (DocumentVersion, error) {
	version, err := s.store.GetVersion(ctx, versionID)
	if err != nil {
		return DocumentVersion{}, err
	}
	if version.DocumentID != documentID {
		return DocumentVersion{}, errors.New("version does not belong to document")
	}

	document, err := s.store.SetCurrentAdoptedVersion(ctx, documentID, versionID)
	if err != nil {
		return DocumentVersion{}, err
	}
	if document.CurrentAdoptedVersionID != versionID {
		return DocumentVersion{}, errors.New("failed to set current adopted version")
	}

	return s.store.GetVersion(ctx, versionID)
}
