package archive

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type MissionArchive struct {
	ID                string `json:"id"`
	MissionID         string `json:"missionId"`
	Status            string `json:"status"`
	ManifestObjectKey string `json:"manifestObjectKey,omitempty"`
	BundleObjectKey   string `json:"bundleObjectKey,omitempty"`
	Hash              string `json:"hash,omitempty"`
}

type Store interface {
	BuildArchiveData(context.Context, string) (ArchiveData, error)
	CreateArchiveRecord(context.Context, string, ArchiveManifest) (MissionArchive, error)
	GetArchiveByMission(context.Context, string) (MissionArchive, error)
}

type Service interface {
	BuildManifest(context.Context, string) (ArchiveManifest, error)
	CreateArchive(context.Context, string) (MissionArchive, ArchiveManifest, error)
	GetArchive(context.Context, string) (MissionArchive, error)
}

type service struct {
	store  Store
	writer ObjectWriter
}

type ObjectWriter interface {
	PutJSON(context.Context, string, any) error
	PutBytes(context.Context, string, []byte, string) error
}

// NewService 创建并返回对应的组件。
func NewService(store Store, writer ...ObjectWriter) Service {
	var selected ObjectWriter
	if len(writer) > 0 {
		selected = writer[0]
	}
	return &service{
		store:  store,
		writer: selected,
	}
}

// BuildManifest builds the requested artifact from the available inputs.
func (s *service) BuildManifest(ctx context.Context, missionID string) (ArchiveManifest, error) {
	data, err := s.store.BuildArchiveData(ctx, missionID)
	if err != nil {
		return ArchiveManifest{}, err
	}
	return BuildManifest(missionID, data), nil
}

// CreateArchive 构建 manifest，写入归档对象，并记录归档行。
func (s *service) CreateArchive(ctx context.Context, missionID string) (MissionArchive, ArchiveManifest, error) {
	// 1. 收集归档输入，并生成 manifest 与 bundle 数据。
	data, err := s.store.BuildArchiveData(ctx, missionID)
	if err != nil {
		return MissionArchive{}, ArchiveManifest{}, err
	}
	manifest := BuildManifest(missionID, data)

	// 2. 在启用对象存储时，持久化 manifest 和 bundle 产物。
	if s.writer != nil {
		if err := s.writer.PutJSON(ctx, archiveManifestObjectKey(missionID), manifest); err != nil {
			return MissionArchive{}, ArchiveManifest{}, err
		}
		bundleBytes, err := buildBundleZip(manifest, BuildBundle(missionID, data))
		if err != nil {
			return MissionArchive{}, ArchiveManifest{}, err
		}
		if err := s.writer.PutBytes(ctx, archiveBundleObjectKey(missionID), bundleBytes, "application/zip"); err != nil {
			return MissionArchive{}, ArchiveManifest{}, err
		}
	}

	// 3. 记录指向对象存储键的归档数据库记录。
	archive, err := s.store.CreateArchiveRecord(ctx, missionID, manifest)
	if err != nil {
		return MissionArchive{}, ArchiveManifest{}, err
	}
	return archive, manifest, nil
}

// GetArchive 返回请求的资源或值。
func (s *service) GetArchive(ctx context.Context, missionID string) (MissionArchive, error) {
	return s.store.GetArchiveByMission(ctx, missionID)
}

type Repository struct {
	queries *sqlc.Queries
}

// NewRepository 创建并返回对应的组件。
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{queries: sqlc.New(pool)}
}

// BuildArchiveData 收集归档所需的文档、决策、审批和运行时引用。
func (r *Repository) BuildArchiveData(ctx context.Context, missionID string) (ArchiveData, error) {
	documentRows, err := r.queries.ListAdoptedDocumentVersionIDsByMission(ctx, missionID)
	if err != nil {
		return ArchiveData{}, err
	}
	decisionIDs, err := r.queries.ListMissionDecisionIDs(ctx, missionID)
	if err != nil {
		return ArchiveData{}, err
	}
	approvalIDs, err := r.queries.ListApprovalIDsByMission(ctx, textValue(missionID))
	if err != nil {
		return ArchiveData{}, err
	}
	runtimeKeyRows, err := r.queries.ListRuntimeSummaryKeysByMission(ctx, missionID)
	if err != nil {
		return ArchiveData{}, err
	}
	documentIDs := make([]string, 0, len(documentRows))
	for _, value := range documentRows {
		if value.Valid {
			documentIDs = append(documentIDs, value.String)
		}
	}
	runtimeKeys := make([]string, 0, len(runtimeKeyRows))
	for _, value := range runtimeKeyRows {
		runtimeKeys = append(runtimeKeys, fmt.Sprint(value))
	}
	return ArchiveData{
		AdoptedDocumentVersionIDs: documentIDs,
		IncludedDecisionIDs:       decisionIDs,
		IncludedApprovalIDs:       approvalIDs,
		RuntimeEventSummaryKeys:   runtimeKeys,
	}, nil
}

// CreateArchiveRecord 创建请求的资源或记录。
func (r *Repository) CreateArchiveRecord(ctx context.Context, missionID string, manifest ArchiveManifest) (MissionArchive, error) {
	row, err := r.queries.CreateMissionArchive(ctx, sqlc.CreateMissionArchiveParams{
		ID:               uuid.NewString(),
		MissionID:        missionID,
		Status:           "uploaded",
		ManifestObjectKey: textValue("archives/" + missionID + "/manifest.json"),
		BundleObjectKey:   textValue("archives/" + missionID + "/bundle.zip"),
		Hash:              textValue(manifest.ManifestHash),
	})
	if err != nil {
		return MissionArchive{}, err
	}
	return archiveFromRow(row), nil
}

// GetArchiveByMission 返回请求的资源或值。
func (r *Repository) GetArchiveByMission(ctx context.Context, missionID string) (MissionArchive, error) {
	row, err := r.queries.GetMissionArchiveByMission(ctx, missionID)
	if err != nil {
		return MissionArchive{}, err
	}
	return archiveFromRow(row), nil
}

// archiveFromRow 实现当前函数行为。
func archiveFromRow(row sqlc.MissionArchive) MissionArchive {
	return MissionArchive{
		ID:                row.ID,
		MissionID:         row.MissionID,
		Status:            row.Status,
		ManifestObjectKey: stringValue(row.ManifestObjectKey),
		BundleObjectKey:   stringValue(row.BundleObjectKey),
		Hash:              stringValue(row.Hash),
	}
}

// archiveManifestObjectKey 实现当前函数行为。
func archiveManifestObjectKey(missionID string) string {
	return "archives/" + missionID + "/manifest.json"
}

// archiveBundleObjectKey 实现当前函数行为。
func archiveBundleObjectKey(missionID string) string {
	return "archives/" + missionID + "/bundle.zip"
}

// buildBundleZip 实现当前函数行为。
func buildBundleZip(manifest ArchiveManifest, bundle ArchiveBundle) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := zip.NewWriter(buffer)

	if err := writeZipJSON(writer, "manifest.json", manifest); err != nil {
		return nil, err
	}
	if err := writeZipJSON(writer, "bundle.json", bundle); err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// writeZipJSON 实现当前函数行为。
func writeZipJSON(writer *zip.Writer, name string, value any) error {
	entry, err := writer.Create(name)
	if err != nil {
		return err
	}
	payload, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	_, err = entry.Write(payload)
	return err
}

// textValue 实现当前函数行为。
func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

// stringValue 实现当前函数行为。
func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

var _ = json.RawMessage{}
