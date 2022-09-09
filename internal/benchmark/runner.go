package benchmark

import (
	"context"
	"sync"
)

type RunnerOption func(r *Runner)

type Runner struct {
	pool  PGPool
	args  *QueryArgsScanner
	group *Group
}

func NewRunner(pool PGPool, args *QueryArgsScanner, opts ...RunnerOption) *Runner {
	r := &Runner{
		pool:  pool,
		args:  args,
		group: NewGroup(pool, 4),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func WithWorkers(count uint) RunnerOption {
	return func(r *Runner) {
		r.group = NewGroup(r.pool, count)
	}
}

func (r *Runner) Run(ctx context.Context, sql string) Summary {
	var (
		wg sync.WaitGroup
		mu sync.Mutex

		summary Summary
	)

	summary.RecordStartTime()

	for r.args.Scan() {
		args := r.args.Data()
		key := args[0]

		wg.Add(1)
		go func() {
			defer wg.Done()
			res := r.group.Benchmark(ctx, key, sql, args)

			mu.Lock()
			defer mu.Unlock()
			summary.Add(res)
		}()
	}

	wg.Wait()
	summary.RecordStopTime()

	return summary
}
