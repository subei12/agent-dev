package authz

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

const (
	CapabilityManageMission    = "manage_mission"
	CapabilityViewTranscripts  = "view_transcripts"
	CapabilityExportTranscript = "export_transcripts"
	CapabilityManageArchive    = "manage_archive"
)

var ErrForbidden = errors.New("forbidden")

type Authorizer interface {
	Require(context.Context, string, string, string) error
}

type Service struct {
	queries *sqlc.Queries
}

func New(db sqlc.DBTX) *Service {
	return &Service{queries: sqlc.New(db)}
}

func (s *Service) Require(ctx context.Context, projectID, actorUserID, capability string) error {
	if actorUserID == "" {
		return ErrForbidden
	}

	grant, err := s.queries.GetProjectAccessGrant(ctx, sqlc.GetProjectAccessGrantParams{
		ProjectID:   projectID,
		ActorUserID: actorUserID,
	})
	if err != nil {
		return ErrForbidden
	}

	switch capability {
	case CapabilityManageMission:
		if grant.CanManageMission {
			return nil
		}
	case CapabilityViewTranscripts:
		if grant.CanViewTranscripts {
			return nil
		}
	case CapabilityExportTranscript:
		if grant.CanExportTranscripts {
			return nil
		}
	case CapabilityManageArchive:
		if grant.CanManageArchive {
			return nil
		}
	default:
		return ErrForbidden
	}

	return ErrForbidden
}

func SeedGrant(ctx context.Context, queries *sqlc.Queries, projectID, actorUserID string) error {
	_, err := queries.UpsertProjectAccessGrant(ctx, sqlc.UpsertProjectAccessGrantParams{
		ID:                  uuid.NewString(),
		ProjectID:           projectID,
		ActorUserID:         actorUserID,
		CanManageMission:    true,
		CanViewTranscripts:  true,
		CanExportTranscripts: false,
		CanManageArchive:    true,
	})
	return err
}
