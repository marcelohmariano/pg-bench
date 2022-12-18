package bench

import (
	"testing"
)

func TestMedianCalculator(t *testing.T) {
	tests := []struct {
		name  string
		input []float64
		want  float64
	}{
		{
			name:  "one number",
			input: []float64{100},
			want:  100,
		},
		{
			name:  "two numbers",
			input: []float64{1, 2},
			want:  1.5,
		},
		{
			name:  "three numbers",
			input: []float64{1, 2, 3},
			want:  2,
		},
		{
			name:  "odd amount of numbers",
			input: []float64{3, 13, 7, 5, 21, 23, 39, 23, 40, 23, 14, 12, 56, 23, 29},
			want:  23,
		},
		{
			name:  "even amount of numbers",
			input: []float64{3, 13, 7, 5, 21, 23, 23, 40, 23, 14, 12, 56, 23, 29},
			want:  22,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m medianCalculator
			got := m.Calculate(tt.input...)
			if got != tt.want {
				t.Errorf("Calculate(): got %v, want %v", got, tt.want)
			}
		})
	}
}
