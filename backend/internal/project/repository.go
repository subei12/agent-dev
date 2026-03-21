package project

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type CreateProjectParams struct {
	Name        string
	Description string
}

type Repository struct {
	queries *sqlc.Queries
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		queries: sqlc.New(pool),
	}
}

func (r *Repository) Create(ctx context.Context, params CreateProjectParams) (sqlc.Project, error) {
	description := pgtype.Text{}
	if params.Description != "" {
		description = pgtype.Text{
			String: params.Description,
			Valid:  true,
		}
	}

	return r.queries.CreateProject(ctx, sqlc.CreateProjectParams{
		ID:          uuid.NewString(),
		Name:        params.Name,
		Description: description,
		Status:      "active",
	})
}
