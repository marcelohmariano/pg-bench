package bench

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5"
)

func TestStats_AddResult(t *testing.T) {
	type result struct {
		duration time.Duration
		err      error
	}

	overallD := 150 * time.Millisecond
	timeElapsedFunc = func(t time.Time) time.Duration { return overallD }

	tests := []struct {
		name    string
		results []result
		want    *Stats
	}{
		{
			name: "successes",
			results: []result{
				{10 * time.Millisecond, nil},
				{20 * time.Millisecond, nil},
				{30 * time.Millisecond, nil},
				{40 * time.Millisecond, nil},
				{50 * time.Millisecond, nil},
			},
			want: &Stats{
				Durations: Durations{
					Min:     10 * time.Millisecond,
					Max:     50 * time.Millisecond,
					Average: 30 * time.Millisecond,
					Median:  30 * time.Millisecond,
					Overall: overallD,
				},
				Statements: Statements{
					Total:     5,
					Succeeded: 5,
					Failed:    0,
				},
			},
		},
		{
			name: "errors",
			results: []result{
				{10 * time.Millisecond, pgx.ErrNoRows},
				{20 * time.Millisecond, pgx.ErrNoRows},
				{30 * time.Millisecond, pgx.ErrNoRows},
				{40 * time.Millisecond, pgx.ErrNoRows},
				{50 * time.Millisecond, pgx.ErrNoRows},
			},
			want: &Stats{
				Durations: Durations{Overall: overallD},
				Statements: Statements{
					Total:     5,
					Succeeded: 0,
					Failed:    5,
				},
			},
		},
		{
			name: "successes and errors",
			results: []result{
				{10 * time.Millisecond, pgx.ErrNoRows},
				{20 * time.Millisecond, nil},
				{30 * time.Millisecond, pgx.ErrNoRows},
				{40 * time.Millisecond, nil},
				{50 * time.Millisecond, pgx.ErrNoRows},
			},
			want: &Stats{
				Durations: Durations{
					Min:     20 * time.Millisecond,
					Max:     40 * time.Millisecond,
					Average: 30 * time.Millisecond,
					Median:  30 * time.Millisecond,
					Overall: overallD,
				},
				Statements: Statements{
					Total:     5,
					Succeeded: 2,
					Failed:    3,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStats()
			for _, r := range tt.results {
				got.AddResult(TaskResult{Duration: r.duration, Error: r.err})
			}
			got.End()

			if diff := got.diff(tt.want); diff != "" {
				t.Errorf("Stats mismatch (-got +want):\n%s", cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestStats_String(t *testing.T) {
	type result struct {
		duration time.Duration
		err      error
	}

	overallD := 150 * time.Millisecond

	tests := []struct {
		name    string
		results []result
		want    string
	}{
		{
			name: "successes and errors",
			results: []result{
				{10 * time.Millisecond, pgx.ErrNoRows},
				{20 * time.Millisecond, nil},
				{30 * time.Millisecond, pgx.ErrNoRows},
				{40 * time.Millisecond, nil},
				{50 * time.Millisecond, pgx.ErrNoRows},
			},
			want: "Statements:\n" +
				"  Total: 5\n" +
				"  Succeeded: 2\n" +
				"  Failed: 3\n\n" +
				"Durations:\n" +
				"  Min: 20ms\n" +
				"  Max: 40ms\n" +
				"  Average: 30ms\n" +
				"  Median: 30ms\n" +
				"  Overall: " + overallD.String(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewStats()
			for _, r := range tt.results {
				res := TaskResult{Duration: r.duration, Error: r.err}
				s.AddResult(res)
			}
			s.End()
			got := s.String()

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("String() mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
