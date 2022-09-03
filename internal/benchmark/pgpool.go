package benchmark

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PGPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

var _ PGPool = (*pgxpool.Pool)(nil)
