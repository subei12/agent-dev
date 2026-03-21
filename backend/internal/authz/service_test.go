package authz

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v4"
)

// TestRequireReturnsForbiddenWhenGrantMissing 验证该路径的预期行为。
func TestRequireReturnsForbiddenWhenGrantMissing(t *testing.T) {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("new mock pool: %v", err)
	}
	defer mock.Close()

	mock.ExpectQuery("select id, project_id").
		WithArgs("proj_1", "demo-user").
		WillReturnError(pgx.ErrNoRows)

	svc := New(mock)
	if err := svc.Require(context.Background(), "proj_1", "demo-user", CapabilityManageMission); err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}
