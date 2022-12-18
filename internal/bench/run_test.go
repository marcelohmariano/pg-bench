package bench

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		ctx     func() context.Context
		input   map[string]time.Duration
		workers int
		wantErr bool
	}{
		{
			name: "success",
			ctx:  func() context.Context { return context.Background() },
			input: map[string]time.Duration{
				"select * from fake_table where id = 1": 100 * time.Millisecond,
				"select * from fake_table where id = 2": 200 * time.Millisecond,
				"select * from fake_table where id = 3": 300 * time.Millisecond,
			},
			workers: 3,
			wantErr: false,
		},
	}

	timeElapsedFunc = func(t time.Time) time.Duration { return 400 * time.Millisecond }

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts := make([]string, 0, len(tt.input))
			for s := range tt.input {
				stmts = append(stmts, s)
			}
			sql := strings.Join(stmts, ";\n") + ";"

			ctx := tt.ctx()
			input := strings.NewReader(sql)
			runner := TaskRunner(func(ctx context.Context, stmt string) TaskResult {
				return TaskResult{Duration: tt.input[stmt]}
			})

			want := NewStats()
			for _, d := range tt.input {
				want.AddResult(TaskResult{Duration: d})
			}
			want.End()

			got, err := Run(ctx, input, runner, NumWorkers(tt.workers))
			if diff := got.diff(want); diff != "" {
				t.Errorf("Run():\n-got +want\n%s", diff)
			}
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}

}
