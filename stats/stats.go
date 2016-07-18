package stats

import (
	"fmt"
	"io"
	"time"
	"sort"
)

type durationsSlice []time.Duration

type Stats struct {
	durations durationsSlice
	len       int
}

func New() *Stats {
	return &Stats{
		durations: make([]time.Duration, 10000000),
	}
}
func (s *Stats) Len() int { return s.len }
func (s *Stats) Swap(i, j int) { s.durations[i], s.durations[j] = s.durations[j], s.durations[i] }
func (s *Stats) Less(i, j int) bool { return s.durations[i] < s.durations[j] }

// Add adds a duration to stats. It will panic with a "index out of range" if
// there is no space.
func (s *Stats) Add(d time.Duration) {
	s.durations[s.len] = d
	s.len++
}

// Print writes textual output of the Stats.
func (s *Stats) Print(suffix string, w io.Writer) {
	if s.len < 2 {
		fmt.Fprint(w, "Histogram (too few values)\n")
	}
	sort.Sort(s)
	min := s.durations[0].Nanoseconds()
	max := s.durations[s.len-1].Nanoseconds()

	sum := int64(0)
	for _, value := range s.durations {
		sum = sum + value.Nanoseconds()
	}
	fmt.Fprintf(w, "Count: %d  Min: %d  Max: %d  Avg: %d\n", s.len, min, max, sum/int64(s.len))
	fmt.Fprint(w, "------------------------------------------------------------\n")

	fmt.Fprintf(w, "pos total ns suffix...\n")
	fmt.Fprintf(w, "%d %d %d %s\n", 0, s.len, min, suffix)
	for i := 1; i < 100; i++ {
		pos := s.len * i / 100
		fmt.Fprintf(w, "%d %d %d %s\n", pos, s.len, s.durations[pos], suffix)
	}
	fmt.Fprintf(w, "%d %d %d %s\n", s.len, s.len, max, suffix)
}
