package stats

import (
	"fmt"
	"io"
	"time"
	"sort"
	"strconv"
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
func (s *Stats) Print(w io.Writer) {
	if s.len < 2 {
		fmt.Fprint(w, "Histogram (too few values)\n")
	}
	sort.Sort(s)
	min := s.durations[0].Nanoseconds()
	max := s.durations[s.len-1].Nanoseconds()

	// Use the largest unit that can represent the minimum time duration.
	unit := time.Nanosecond
	for _, u := range []time.Duration{time.Microsecond, time.Millisecond, time.Second} {
		if min <= u.Nanoseconds() {
			break
		}
		unit = u
	}
	unitFactor := unit.Nanoseconds()

	fmt.Fprintf(w, "Histogram (unit: %v, %dns)\n", unit, unit.Nanoseconds())
	fmt.Fprintf(w, "Count: %d  Min: %d  Max: %d  Avg: \n", s.len, min/unitFactor, max/unitFactor)
	fmt.Fprint(w, "------------------------------------------------------------\n")

	maxBucketDigitLen := len(strconv.FormatInt(int64(max)/unitFactor, 10))
	if maxBucketDigitLen < 3 {
		maxBucketDigitLen = 3
	}
	maxCountDigitLen := len(strconv.FormatInt(int64(s.len), 10))
	fmt.Fprintf(w, "(%*s, %*d) %*d %*d %3d%%\n", maxBucketDigitLen, "-inf", maxBucketDigitLen, int64(min)/unitFactor, maxCountDigitLen, 0, maxCountDigitLen, 0, 0)
	lastDuration := int64(min)
	lastCount := 0
	for i := 1; i < 100; i++ {
		pos := s.len * i / 100
		currentDuration := int64(s.durations[pos].Nanoseconds())
		if currentDuration < lastDuration {
			continue
		}
		// Skip all the values equal with the current one because we
		// promise a ')'.
		for pos+1 < s.len && s.durations[pos+1] == s.durations[pos] {
			pos++
		}
		currentDuration = int64(s.durations[pos].Nanoseconds())
		fmt.Fprintf(w, "[%*d, %*d) %*d %*d %3d%%\n", maxBucketDigitLen, lastDuration/unitFactor, maxBucketDigitLen, currentDuration/unitFactor, maxCountDigitLen, pos - lastCount, maxCountDigitLen, pos, i)
		lastDuration = currentDuration
		lastCount = pos
	}
	fmt.Fprintf(w, "[%*d, %*d] %*d %*d %3d%%\n", maxBucketDigitLen, lastDuration/unitFactor, maxBucketDigitLen, int64(max)/unitFactor, maxCountDigitLen, s.len - lastCount, maxCountDigitLen, s.len, 100)
}
