//go:build integration

package project

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TestCreateProject 验证该路径的预期行为。
func TestCreateProject(t *testing.T) {
	ctx := context.Background()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable"
	}

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	defer pool.Close()

	repo := NewRepository(pool)
	project, err := repo.Create(ctx, CreateProjectParams{
		Name: "Demo",
	})
	if err != nil {
		t.Fatalf("create project: %v", err)
	}

	if project.Name != "Demo" {
		t.Fatalf("expected project name Demo, got %q", project.Name)
	}
}
