package benchmark

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stretchr/testify/assert"
)

type poolStub struct{}

func (p poolStub) Acquire(ctx context.Context) (*pgxpool.Conn, error) {
	return nil, nil
}

func (p poolStub) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return nil, nil
}

func TestGroup_WorkerOf(t *testing.T) {
	tests := []struct {
		name    string
		workers uint
		keys    []any
	}{
		{
			name:    "with 1 worker",
			workers: 1,
			keys: []any{"host_000000", "host_000001", "host_000002", "host_000003", "host_000004",
				"host_000005", "host_000006", "host_000007", "host_000008", "host_000009",
				"host_000010",
			},
		},
		{
			name:    "with 2 workers",
			workers: 2,
			keys: []any{"host_000000", "host_000001", "host_000002", "host_000003", "host_000004",
				"host_000005", "host_000006", "host_000007", "host_000008", "host_000009",
				"host_000010",
			},
		},
		{
			name:    "with 4 workers",
			workers: 4,
			keys: []any{"host_000000", "host_000001", "host_000002", "host_000003", "host_000004",
				"host_000005", "host_000006", "host_000007", "host_000008", "host_000009",
				"host_000010",
			},
		},
		{
			name:    "with 10 workers",
			workers: 10,
			keys: []any{"host_000000", "host_000001", "host_000002", "host_000003", "host_000004",
				"host_000005", "host_000006", "host_000007", "host_000008", "host_000009",
				"host_000010",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := NewGroup(poolStub{}, tt.workers)
			for _, key := range tt.keys {
				worker1, err := group.WorkerOf(key)
				assert.NotNil(t, worker1)
				assert.Nil(t, err)

				worker2, err := group.WorkerOf(key)
				assert.NotNil(t, worker2)
				assert.Nil(t, err)

				assert.Equal(t, worker1.ID, worker2.ID)
			}
		})
	}
}
