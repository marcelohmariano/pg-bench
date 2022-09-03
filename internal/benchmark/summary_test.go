package benchmark

import (
	"bufio"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func openExecTimesFile(t *testing.T, name string) (*bufio.Scanner, func()) {
	t.Helper()

	f, err := os.Open(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	return scanner, func() { _ = f.Close() }
}

func TestSummary(t *testing.T) {
	tests := []struct {
		name                      string
		fixture                   string
		expectedNumberOfQueries   int64
		expectedNumberOfSuccesses int64
		expectedNumberOfErrors    int64
		expectedMinQueryTime      time.Duration
		expectedMaxQueryTime      time.Duration
		expectedMedianQueryTime   time.Duration
		expectedAvgQueryTime      time.Duration
		expectedOverallQueryTime  string
	}{
		{
			name:                      "with sample exec times",
			fixture:                   "exec_times.txt",
			expectedNumberOfQueries:   200,
			expectedNumberOfSuccesses: 200,
			expectedNumberOfErrors:    0,
			expectedMinQueryTime:      9820166,
			expectedMaxQueryTime:      204942208,
			expectedMedianQueryTime:   19078437,
			expectedAvgQueryTime:      26050722,
			expectedOverallQueryTime:  "100ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner, tearDown := openExecTimesFile(t, tt.fixture)
			defer tearDown()

			expectedOverallQueryTime, _ := time.ParseDuration(tt.expectedOverallQueryTime)
			timeSinceFunc = func(t time.Time) time.Duration {
				return expectedOverallQueryTime
			}

			var stats Summary
			stats.Start()

			for scanner.Scan() {
				nsec, err := strconv.ParseInt(scanner.Text(), 10, 64)
				if err != nil {
					t.Fatal(err)
				}
				result := Result{ExecTime: time.Duration(nsec)}
				stats.Add(result)
			}

			stats.Done()

			assert.Equal(t, tt.expectedNumberOfQueries, stats.NumberOfQueries())
			assert.Equal(t, tt.expectedNumberOfSuccesses, stats.NumberOfSuccesses())
			assert.Equal(t, tt.expectedNumberOfErrors, stats.NumberOfErrors())
			assert.Equal(t, tt.expectedMinQueryTime, stats.MinQueryTime())
			assert.Equal(t, tt.expectedMaxQueryTime, stats.MaxQueryTime())
			assert.Equal(t, tt.expectedMedianQueryTime, stats.MedianQueryTime())
			assert.Equal(t, tt.expectedAvgQueryTime, stats.AvgQueryTime())
			assert.Equal(t, expectedOverallQueryTime, stats.OverallQueryTime())

			stats.Start()

			assert.Equal(t, int64(0), stats.NumberOfQueries())
			assert.Equal(t, int64(0), stats.NumberOfSuccesses())
			assert.Equal(t, int64(0), stats.NumberOfErrors())
			assert.Equal(t, time.Duration(0), stats.MinQueryTime())
			assert.Equal(t, time.Duration(0), stats.MaxQueryTime())
			assert.Equal(t, time.Duration(0), stats.MedianQueryTime())
			assert.Equal(t, time.Duration(0), stats.AvgQueryTime())
			assert.Equal(t, time.Duration(0), stats.OverallQueryTime())
		})
	}
}
