package bench

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Task = string

type TaskResult struct {
	Error    error
	Duration time.Duration
}

type TaskRunner func(ctx context.Context, task Task) TaskResult

func NewTaskRunner(pool *pgxpool.Pool) TaskRunner {
	return func(ctx context.Context, task Task) TaskResult {
		var (
			start   time.Time
			elapsed time.Duration
		)

		err := pool.AcquireFunc(ctx, func(conn *pgxpool.Conn) error {
			start = time.Now()
			_, err := conn.Exec(ctx, task)
			elapsed = time.Since(start)
			return err
		})

		return TaskResult{Duration: elapsed, Error: err}
	}
}

type taskProducer struct {
	scanner *bufio.Scanner
	C       chan Task
}

func newTaskProducer(input io.Reader) *taskProducer {
	p := &taskProducer{
		scanner: bufio.NewScanner(input),
		C:       make(chan Task),
	}
	p.scanner.Split(p.nextStmt)
	return p
}

func (p *taskProducer) Run(ctx context.Context) error {
	for p.scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case p.C <- p.scanner.Text():
		}
	}
	return nil
}

func (p *taskProducer) Close() {
	close(p.C)
}

func (p *taskProducer) nextStmt(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	c := 0
	for {
		i := bytes.IndexByte(data, ';')
		if i < 0 {
			break
		}

		t := data[0:i]

		token = bytes.TrimSpace(t)
		if len(token) == 0 {
			data, c = data[1:], c+1
			continue
		}

		return i + c + 1, token, nil
	}

	if atEOF {
		return len(data), bytes.TrimSpace(data), nil
	}

	return 0, nil, nil
}

type taskRunners struct {
	runner   TaskRunner
	producer *taskProducer
	handler  *taskResultHandler
	wg       sync.WaitGroup
	workers  int
}

func newTaskRunners(r TaskRunner, p *taskProducer, h *taskResultHandler, workers int) taskRunners {
	return taskRunners{
		runner:   r,
		producer: p,
		handler:  h,
		workers:  workers,
	}
}

func (r *taskRunners) Run(ctx context.Context) {
	r.wg.Add(r.workers)
	for i := 0; i < r.workers; i++ {
		go func() {
			defer r.wg.Done()
			r.run(ctx)
		}()
	}
}

func (r *taskRunners) Wait() {
	r.wg.Wait()
}

func (r *taskRunners) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case s, ok := <-r.producer.C:
			if !ok {
				return
			}
			r.handler.C <- r.runner(ctx, s)
		}
	}
}

type taskResultHandler struct {
	stats *Stats
	errc  chan error
	C     chan TaskResult
}

func newTaskResultHandler(s *Stats) *taskResultHandler {
	return &taskResultHandler{
		stats: s,
		errc:  make(chan error, 1),
		C:     make(chan TaskResult),
	}
}

func (h *taskResultHandler) Close() {
	close(h.C)
}

func (h *taskResultHandler) Err() error {
	return <-h.errc
}

func (h *taskResultHandler) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			h.errc <- ctx.Err()
		case res, ok := <-h.C:
			if !ok {
				h.errc <- nil
				return
			}
			h.stats.AddResult(res)
		}
	}
}
