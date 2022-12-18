package bench

import (
	"context"
	"io"
	"runtime"
)

type options struct {
	workers int
}

func newOptions(opts []Option) options {
	o := options{workers: runtime.NumCPU()}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

type Option func(o *options)

func NumWorkers(n int) Option {
	return func(o *options) {
		o.workers = n
	}
}

func Run(ctx context.Context, input io.Reader, runner TaskRunner, opts ...Option) (*Stats, error) {
	stats := NewStats()
	defer stats.End()

	handler := newTaskResultHandler(stats)
	producer := newTaskProducer(input)

	o := newOptions(opts)
	runners := newTaskRunners(runner, producer, handler, o.workers)

	go handler.Run(ctx)
	go runners.Run(ctx)

	if err := producer.Run(ctx); err != nil {
		return nil, err
	}

	producer.Close()
	runners.Wait()
	handler.Close()

	if err := handler.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
