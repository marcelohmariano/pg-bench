package benchmark

import "context"

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

func (r Runner) Run(ctx context.Context, sql string) (*Summary, error) {
	for r.args.Scan() {
		args := r.args.Data()
		key := args[0]

		task := NewTask(key, sql, args)

		if err := r.group.AddTask(ctx, task); err != nil {
			return nil, err
		}
	}

	return r.group.Wait(), nil
}
