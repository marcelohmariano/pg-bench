package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMedian(t *testing.T) {
	tests := []struct {
		name     string
		input    []float64
		expected float64
	}{
		{
			name:     "with 1 number",
			input:    []float64{100},
			expected: 100,
		},
		{
			name:     "with 2 numbers",
			input:    []float64{1, 2},
			expected: 1.5,
		},
		{
			name:     "with 3 numbers",
			input:    []float64{1, 2, 3},
			expected: 2,
		},
		{
			name:     "with odd amount of numbers",
			input:    []float64{3, 13, 7, 5, 21, 23, 39, 23, 40, 23, 14, 12, 56, 23, 29},
			expected: 23,
		},
		{
			name:     "with even amount of numbers",
			input:    []float64{3, 13, 7, 5, 21, 23, 23, 40, 23, 14, 12, 56, 23, 29},
			expected: 22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var median Median
			median.Add(tt.input...)
			assert.Equal(t, tt.expected, median.Value())
		})
	}
}
