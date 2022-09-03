package benchmark

import (
	"math"
	"time"
)

var timeSinceFunc = func(t time.Time) time.Duration {
	return time.Since(t)
}

type Result struct {
	ExecTime time.Duration
	Err      error
}

type stats struct {
	median Median

	Queries   int64
	Successes int64
	Errors    int64

	MinTime    time.Duration
	MaxTime    time.Duration
	AvgTime    float64
	MedianTime float64

	StartTime   time.Time
	OverallTime time.Duration
}

func newStats() *stats {
	return &stats{
		MinTime:   math.MaxInt64,
		StartTime: time.Now(),
	}
}

func (s *stats) IncQueries() {
	s.Queries++
}

func (s *stats) IncSuccesses() {
	s.Successes++
}

func (s *stats) IncErrors() {
	s.Errors++
}

func (s *stats) UpdateMinTime(d time.Duration) {
	if d < s.MinTime {
		s.MinTime = d
	}
}

func (s *stats) UpdateMaxTime(d time.Duration) {
	if d > s.MaxTime {
		s.MaxTime = d
	}
}

func (s *stats) UpdateMedianTime(d time.Duration) {
	s.median.Add(float64(d))
	s.MedianTime = s.median.Value()
}

func (s *stats) UpdateAvgTime(d time.Duration) {
	successes := float64(s.Successes - 1)
	oldSum := s.AvgTime * successes
	s.AvgTime = (oldSum + float64(d)) / (successes + 1)
}

func (s *stats) UpdateOverallTime() {
	s.OverallTime = timeSinceFunc(s.StartTime)
}

type Summary struct {
	stats *stats
}

func (s *Summary) Start() {
	s.stats = newStats()
}

func (s *Summary) Done() {
	s.stats.UpdateOverallTime()
}

func (s *Summary) Add(result Result) {
	s.stats.IncQueries()

	if result.Err != nil {
		s.stats.IncErrors()
		return
	}

	s.stats.IncSuccesses()
	execTime := result.ExecTime

	s.stats.UpdateMinTime(execTime)
	s.stats.UpdateMaxTime(execTime)
	s.stats.UpdateAvgTime(execTime)
	s.stats.UpdateMedianTime(execTime)
}

func (s *Summary) NumberOfQueries() int64 {
	return s.stats.Queries
}

func (s *Summary) NumberOfSuccesses() int64 {
	return s.stats.Successes
}

func (s *Summary) NumberOfErrors() int64 {
	return s.stats.Errors
}

func (s *Summary) MinQueryTime() time.Duration {
	if s.stats.Queries != 0 {
		return s.stats.MinTime
	}
	return 0
}

func (s *Summary) MaxQueryTime() time.Duration {
	return s.stats.MaxTime
}

func (s *Summary) AvgQueryTime() time.Duration {
	return time.Duration(s.stats.AvgTime)
}

func (s *Summary) MedianQueryTime() time.Duration {
	return time.Duration(s.stats.MedianTime)
}

func (s *Summary) OverallQueryTime() time.Duration {
	return s.stats.OverallTime
}
