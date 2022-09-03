package benchmark

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueryArgsScanner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		comma    rune
		header   bool
		lines    int
		expected []QueryArgs
	}{
		{
			name:   "skipping header",
			input:  "hostname,start_time,end_time\nhost_000008,2017-01-01 08:59:22,2017-01-01 09:59:22",
			comma:  ',',
			header: true,
			lines:  1,
			expected: []QueryArgs{
				{
					"host_000008",
					"2017-01-01 08:59:22",
					"2017-01-01 09:59:22",
				},
			},
		},
		{
			name:   "not skipping header",
			input:  "host_000008,2017-01-01 08:59:22,2017-01-01 09:59:22",
			comma:  ',',
			header: false,
			lines:  1,
			expected: []QueryArgs{
				{
					"host_000008",
					"2017-01-01 08:59:22",
					"2017-01-01 09:59:22",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := NewQueryArgsScanner(
				WithReader(strings.NewReader(tt.input)),
				WithCSVDelim(tt.comma),
				WithCSVHeader(tt.header),
			)
			defer func() { _ = args.Close() }()

			for i := 0; i < tt.lines; i++ {
				assert.True(t, args.Scan())
				assert.Nil(t, args.Err())
				assert.Equal(t, tt.expected[i], args.Data())

				assert.False(t, args.Scan())
				assert.ErrorIs(t, args.Err(), io.EOF)
				assert.Nil(t, args.Data())
			}
		})
	}
}
