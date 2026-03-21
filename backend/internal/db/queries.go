package db

import (
	"github.com/jackc/pgx/v5"

	"github.com/your-org/agent-platform/internal/db/sqlc"
)

type DBTX interface {
	pgx.Tx
}

// NewQueries 创建并返回对应的组件。
func NewQueries(db pgx.Tx) *sqlc.Queries {
	return sqlc.New(db)
}
