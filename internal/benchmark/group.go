package benchmark

import (
	"bytes"
	"context"
	"encoding/gob"
	"hash/maphash"
)

var hashSeed = maphash.MakeSeed()

type Group struct {
	pool    PGPool
	workers []*Worker
}

func NewGroup(pool PGPool, workers uint) *Group {
	g := &Group{pool: pool}
	g.workers = make([]*Worker, workers)

	for i := 0; i < int(workers); i++ {
		g.workers[i] = NewWorker(pool, WorkerID(i+1))
	}

	return g
}

func (g *Group) Benchmark(ctx context.Context, key any, sql string, args QueryArgs) Result {
	worker, err := g.WorkerOf(key)
	if err != nil {
		return Result{Err: err}
	}
	return worker.Benchmark(ctx, sql, args)
}

func (g *Group) WorkerOf(key any) (*Worker, error) {
	hash, err := g.hash(key)
	if err != nil {
		return nil, err
	}

	idx := hash % uint64(len(g.workers))
	return g.workers[idx], nil
}

func (g *Group) hash(key any) (uint64, error) {
	var h maphash.Hash
	h.SetSeed(hashSeed)

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(key); err != nil {
		return 0, err
	}

	_, err := h.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}

	return h.Sum64(), nil
}
