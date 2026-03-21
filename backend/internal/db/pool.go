package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPool 创建并返回对应的组件。
func NewPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
	return pgxpool.New(ctx, databaseURL)
}
