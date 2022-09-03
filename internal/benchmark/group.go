package benchmark

import (
	"bytes"
	"context"
	"encoding/gob"
	"hash/maphash"
	"sync"
)

type Task struct {
	Key  any
	SQL  string
	Args QueryArgs
}

func NewTask(key any, sql string, args QueryArgs) *Task {
	return &Task{
		Key:  key,
		SQL:  sql,
		Args: args,
	}
}

type Group struct {
	mu     sync.Mutex
	wg     sync.WaitGroup
	once   *sync.Once
	hasher maphash.Hash

	pool    PGPool
	workers []*Worker

	summary *Summary
}

func NewGroup(pool PGPool, workers uint) *Group {
	return &Group{
		pool:    pool,
		once:    new(sync.Once),
		summary: new(Summary),
		workers: make([]*Worker, workers),
	}
}

func (g *Group) AddTask(ctx context.Context, task *Task) error {
	worker, err := g.WorkerOf(task.Key)
	if err != nil {
		return err
	}
	return g.runTask(ctx, task, worker)
}

func (g *Group) runTask(ctx context.Context, task *Task, worker *Worker) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	g.once.Do(func() { g.summary.Start() })

	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		res := worker.Benchmark(ctx, task.SQL, task.Args)

		g.mu.Lock()
		defer g.mu.Unlock()
		g.summary.Add(res)
	}()

	return nil
}

func (g *Group) Wait() *Summary {
	g.wg.Wait()
	g.summary.Done()
	g.once = new(sync.Once)
	return g.summary
}

func (g *Group) WorkerOf(key any) (*Worker, error) {
	hash, err := g.hash(key)
	if err != nil {
		return nil, err
	}

	idx := hash % uint64(len(g.workers))

	if g.workers[idx] == nil {
		g.workers[idx] = NewWorker(g.pool, WorkerID(idx+1))
	}

	return g.workers[idx], nil
}

func (g *Group) hash(key any) (uint64, error) {
	g.hasher.Reset()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(key); err != nil {
		return 0, err
	}

	_, err := g.hasher.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}

	return g.hasher.Sum64(), nil
}
