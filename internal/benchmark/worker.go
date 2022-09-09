package benchmark

import (
	"context"
	"time"
)

type WorkerID uint64

type Worker struct {
	pool PGPool

	ID WorkerID
}

func NewWorker(pool PGPool, id WorkerID) *Worker {
	return &Worker{pool: pool, ID: id}
}

func (w *Worker) Benchmark(ctx context.Context, sql string, args QueryArgs) Result {
	conn, err := w.pool.Acquire(ctx)
	if err != nil {
		return Result{Err: err}
	}
	defer conn.Release()

	start := time.Now()
	rows, err := conn.Query(ctx, sql, args...)
	elapsed := time.Since(start)

	if err != nil {
		return Result{Err: err}
	}
	defer rows.Close()

	return Result{Duration: elapsed}
}
