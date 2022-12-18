package bench

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

var timeElapsedFunc = time.Since

type Statements struct {
	Total     int
	Succeeded int
	Failed    int
}

func (s *Statements) update(res TaskResult) (done bool) {
	s.Total++
	if res.Error != nil {
		s.Failed++
		return true
	}
	s.Succeeded++
	return false
}

func (s *Statements) String() string {
	var b strings.Builder

	_, _ = fmt.Fprintln(&b, "Statements:")
	_, _ = fmt.Fprintf(&b, "  Total: %d\n", s.Total)
	_, _ = fmt.Fprintf(&b, "  Succeeded: %d\n", s.Succeeded)
	_, _ = fmt.Fprintf(&b, "  Failed: %d", s.Failed)

	return b.String()
}

type Durations struct {
	mc medianCalculator

	Min     time.Duration
	Max     time.Duration
	Average time.Duration
	Median  time.Duration
	Overall time.Duration
}

func (ds *Durations) update(res TaskResult, totalStmts int) {
	if totalStmts == 0 {
		return
	}

	d := res.Duration

	if d < ds.Min || ds.Min == 0 {
		ds.Min = d
	}
	if d > ds.Max {
		ds.Max = d
	}

	oldSum := float64(ds.Average) * float64(totalStmts-1)
	ds.Average = time.Duration((oldSum + float64(d)) / float64(totalStmts))
	ds.Median = time.Duration(ds.mc.Calculate(float64(d)))
}

func (ds *Durations) String() string {
	var b strings.Builder

	_, _ = fmt.Fprintln(&b, "Durations:")

	empty := ds.Min == 0 && ds.Max == 0 && ds.Average == 0 && ds.Median == 0
	if !empty {
		_, _ = fmt.Fprintf(&b, "  Min: %v\n", ds.Min)
		_, _ = fmt.Fprintf(&b, "  Max: %v\n", ds.Max)
		_, _ = fmt.Fprintf(&b, "  Average: %v\n", ds.Average)
		_, _ = fmt.Fprintf(&b, "  Median: %v\n", ds.Median)
	}
	_, _ = fmt.Fprintf(&b, "  Overall: %v", ds.Overall)

	return b.String()
}

type Stats struct {
	start time.Time

	Durations  Durations
	Statements Statements
}

func NewStats() *Stats {
	return &Stats{start: time.Now()}
}

func (s *Stats) End() {
	s.Durations.Overall = timeElapsedFunc(s.start)
}

func (s *Stats) AddResult(res TaskResult) {
	if done := s.Statements.update(res); done {
		return
	}
	s.Durations.update(res, s.Statements.Succeeded)
}

func (s *Stats) String() string {
	var b strings.Builder
	b.WriteString(s.Statements.String())
	b.WriteString("\n\n")
	b.WriteString(s.Durations.String())
	return b.String()
}

func (s *Stats) diff(o *Stats) string {
	return cmp.Diff(s, o, cmpopts.IgnoreUnexported(Stats{}, Durations{}, Statements{}))
}
