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

func (s *service) BuildManifest(ctx context.Context, missionID string) (ArchiveManifest, error) {
	data, err := s.store.BuildArchiveData(ctx, missionID)
	if err != nil {
		return ArchiveManifest{}, err
	}
	return BuildManifest(missionID, data), nil
}

func (s *service) CreateArchive(ctx context.Context, missionID string) (MissionArchive, ArchiveManifest, error) {
	data, err := s.store.BuildArchiveData(ctx, missionID)
	if err != nil {
		return MissionArchive{}, ArchiveManifest{}, err
	}
	manifest := BuildManifest(missionID, data)
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
	archive, err := s.store.CreateArchiveRecord(ctx, missionID, manifest)
	if err != nil {
		return MissionArchive{}, ArchiveManifest{}, err
	}
	return archive, manifest, nil
}

func (s *service) GetArchive(ctx context.Context, missionID string) (MissionArchive, error) {
	return s.store.GetArchiveByMission(ctx, missionID)
}

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{queries: sqlc.New(pool)}
}

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

func (r *Repository) GetArchiveByMission(ctx context.Context, missionID string) (MissionArchive, error) {
	row, err := r.queries.GetMissionArchiveByMission(ctx, missionID)
	if err != nil {
		return MissionArchive{}, err
	}
	return archiveFromRow(row), nil
}

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

func archiveManifestObjectKey(missionID string) string {
	return "archives/" + missionID + "/manifest.json"
}

func archiveBundleObjectKey(missionID string) string {
	return "archives/" + missionID + "/bundle.zip"
}

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

func textValue(value string) pgtype.Text {
	if value == "" {
		return pgtype.Text{}
	}
	return pgtype.Text{String: value, Valid: true}
}

func stringValue(value pgtype.Text) string {
	if !value.Valid {
		return ""
	}
	return value.String
}

var _ = json.RawMessage{}
