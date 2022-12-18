package bench

import (
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
)

func TestTaskProducer(t *testing.T) {
	tests := []struct {
		name      string
		semicolon bool
		input     []string
		want      []string
	}{
		{
			name: "no empty statements",
			input: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
		{
			name:      "no empty statements with final semicolon",
			semicolon: true,
			input: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
		{
			name: "empty statements",
			input: []string{
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 1",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 2",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 3",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 4",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 5",
				";", ";", ";", ";", ";",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
		{
			name:      "empty statements with final semicolon",
			semicolon: true,
			input: []string{
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 1",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 2",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 3",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 4",
				";", ";", ";", ";", ";",
				"select * from fake_table where id = 5",
				";", ";", ";", ";", ";",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
		{
			name: "empty lines",
			input: []string{
				"\t", "\t", "\t", "\t", "\t",
				"select * from fake_table where id = 1",
				"", "", "", "", "",
				"select * from fake_table where id = 2",
				"\n", "\n", "\n", "\n", "\n",
				"select * from fake_table where id = 3",
				"", "", "", "", "",
				"select * from fake_table where id = 4",
				"\n", "\n", "\n", "\n", "\n",
				"select * from fake_table where id = 5",
				"", "", "", "", "",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
		{
			name:      "empty lines with final semicolon",
			semicolon: true,
			input: []string{
				"\t", "\t", "\t", "\t", "\t",
				"select * from fake_table where id = 1",
				"", "", "", "", "",
				"select * from fake_table where id = 2",
				"\n", "\n", "\n", "\n", "\n",
				"select * from fake_table where id = 3",
				"", "", "", "", "",
				"select * from fake_table where id = 4",
				"\n", "\n", "\n", "\n", "\n",
				"select * from fake_table where id = 5",
				"", "", "", "", "",
			},
			want: []string{
				"select * from fake_table where id = 1",
				"select * from fake_table where id = 2",
				"select * from fake_table where id = 3",
				"select * from fake_table where id = 4",
				"select * from fake_table where id = 5",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql := strings.Join(tt.input, ";\n")
			if tt.semicolon {
				sql += ";"
			}

			p := newTaskProducer(strings.NewReader(sql))
			go func() {
				_ = p.Run(context.TODO())
			}()

			got := make([]string, 0, len(tt.want))
			for range tt.want {
				got = append(got, <-p.C)
			}

			select {
			case s := <-p.C:
				t.Errorf("<-taskProducer.C: got %s, want %s", s, "")
			default:
			}

			p.Close()

			var ok bool
			select {
			case _, ok = <-p.C:
				assert.False(t, ok, "<-taskProducer.C: expected close chan")
			default:
				t.Errorf("<-taskProducer.C: expected close chan")
			}

			if diff := cmp.Diff(got, tt.want); diff != "" {
				t.Errorf("Statements mismatch (-got +want):\n%s", diff)
			}
		})
	}
}
