package document

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type Store interface {
	CreateDocument(context.Context, CreateDocumentCmd) (Document, error)
	GetDocument(context.Context, string, string) (Document, error)
	ListDocuments(context.Context, string) ([]Document, error)
	NextVersionNumber(context.Context, string) (int32, error)
	CreateVersion(context.Context, CreateVersionCmd, int32) (DocumentVersion, error)
	GetVersion(context.Context, string) (DocumentVersion, error)
	ListVersions(context.Context, string) ([]DocumentVersion, error)
	UpdateVersionStatus(context.Context, string, string) (DocumentVersion, error)
	ClearAdoptedVersions(context.Context, string) error
	SetCurrentAdoptedVersion(context.Context, string, string) (Document, error)
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateDocument(ctx context.Context, cmd CreateDocumentCmd) (Document, error) {
	row, err := sqlc.New(r.pool).CreateSharedDocument(ctx, sqlc.CreateSharedDocumentParams{
		ID:                      uuid.NewString(),
		MissionID:               cmd.MissionID,
		Kind:                    cmd.Kind,
		Title:                   cmd.Title,
		CurrentAdoptedVersionID: pgtype.Text{},
	})
	if err != nil {
		return Document{}, err
	}
	return documentFromRow(row), nil
}

func (r *Repository) GetDocument(ctx context.Context, missionID, documentID string) (Document, error) {
	row, err := sqlc.New(r.pool).GetSharedDocument(ctx, sqlc.GetSharedDocumentParams{
		ID:        documentID,
		MissionID: missionID,
	})
	if err != nil {
		return Document{}, err
	}
	return documentFromRow(row), nil
}

func (r *Repository) ListDocuments(ctx context.Context, missionID string) ([]Document, error) {
	rows, err := sqlc.New(r.pool).ListSharedDocumentsByMission(ctx, missionID)
	if err != nil {
		return nil, err
	}
	documents := make([]Document, 0, len(rows))
	for _, row := range rows {
		documents = append(documents, documentFromRow(row))
	}
	return documents, nil
}

func (r *Repository) NextVersionNumber(ctx context.Context, documentID string) (int32, error) {
	return sqlc.New(r.pool).GetNextSharedDocumentVersionNumber(ctx, documentID)
}

func (r *Repository) CreateVersion(ctx context.Context, cmd CreateVersionCmd, version int32) (DocumentVersion, error) {
	row, err := sqlc.New(r.pool).CreateSharedDocumentVersion(ctx, sqlc.CreateSharedDocumentVersionParams{
		ID:                uuid.NewString(),
		DocumentID:        cmd.DocumentID,
		Version:           version,
		ParentVersionID:   textValue(cmd.ParentVersionID),
		Status:            "proposed",
		ContentFormat:     cmd.ContentFormat,
		StorageKind:       cmd.StorageKind,
		ObjectKey:         textValue(cmd.ObjectKey),
		ContentHash:       cmd.ContentHash,
		ContentText:       cmd.ContentText,
		ProducedByAgentID: textValue(cmd.ProducedByAgentID),
		ProducedByRunID:   textValue(cmd.ProducedByRunID),
		SourceRoundID:     textValue(cmd.SourceRoundID),
	})
	if err != nil {
		return DocumentVersion{}, err
	}
	return documentVersionFromRow(row), nil
}

func (r *Repository) GetVersion(ctx context.Context, versionID string) (DocumentVersion, error) {
	row, err := sqlc.New(r.pool).GetSharedDocumentVersion(ctx, versionID)
	if err != nil {
		return DocumentVersion{}, err
	}
	return documentVersionFromRow(row), nil
}

func (r *Repository) ListVersions(ctx context.Context, documentID string) ([]DocumentVersion, error) {
	rows, err := sqlc.New(r.pool).ListSharedDocumentVersions(ctx, documentID)
	if err != nil {
		return nil, err
	}
	versions := make([]DocumentVersion, 0, len(rows))
	for _, row := range rows {
		versions = append(versions, documentVersionFromRow(row))
	}
	return versions, nil
}

func (r *Repository) UpdateVersionStatus(ctx context.Context, versionID, status string) (DocumentVersion, error) {
	row, err := sqlc.New(r.pool).UpdateSharedDocumentVersionStatus(ctx, sqlc.UpdateSharedDocumentVersionStatusParams{
		ID:     versionID,
		Status: status,
	})
	if err != nil {
		return DocumentVersion{}, err
	}
	return documentVersionFromRow(row), nil
}

func (r *Repository) ClearAdoptedVersions(ctx context.Context, documentID string) error {
	return sqlc.New(r.pool).ClearAdoptedSharedDocumentVersions(ctx, documentID)
}

func (r *Repository) SetCurrentAdoptedVersion(ctx context.Context, documentID, versionID string) (Document, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return Document{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	q := sqlc.New(tx)
	if err := q.ClearAdoptedSharedDocumentVersions(ctx, documentID); err != nil {
		return Document{}, err
	}
	if _, err := q.UpdateSharedDocumentVersionStatus(ctx, sqlc.UpdateSharedDocumentVersionStatusParams{
		ID:     versionID,
		Status: "adopted",
	}); err != nil {
		return Document{}, err
	}

	row, err := q.SetSharedDocumentCurrentAdoptedVersion(ctx, sqlc.SetSharedDocumentCurrentAdoptedVersionParams{
		ID:                      documentID,
		CurrentAdoptedVersionID: textValue(versionID),
	})
	if err != nil {
		return Document{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Document{}, err
	}
	return documentFromRow(row), nil
}

func documentFromRow(row sqlc.SharedDocument) Document {
	return Document{
		ID:                      row.ID,
		MissionID:               row.MissionID,
		Kind:                    row.Kind,
		Title:                   row.Title,
		CurrentAdoptedVersionID: stringValue(row.CurrentAdoptedVersionID),
		CreatedAt:               row.CreatedAt.Time,
		UpdatedAt:               row.UpdatedAt.Time,
	}
}

func documentVersionFromRow(row sqlc.SharedDocumentVersion) DocumentVersion {
	return DocumentVersion{
		ID:                row.ID,
		DocumentID:        row.DocumentID,
		Version:           row.Version,
		ParentVersionID:   stringValue(row.ParentVersionID),
		Status:            row.Status,
		ContentFormat:     row.ContentFormat,
		StorageKind:       row.StorageKind,
		ObjectKey:         stringValue(row.ObjectKey),
		ContentHash:       row.ContentHash,
		ContentText:       row.ContentText,
		ProducedByAgentID: stringValue(row.ProducedByAgentID),
		ProducedByRunID:   stringValue(row.ProducedByRunID),
		SourceRoundID:     stringValue(row.SourceRoundID),
		CreatedAt:         row.CreatedAt.Time,
	}
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

var _ pgx.Tx = nil
