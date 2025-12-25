package dates

import (
	"strconv"
	"time"
)

const (
	FmtDate              = "2006-01-02"
	FmtTime              = "15:04:05"
	FmtDateTime          = "2006-01-02 15:04:05"
	FmtDateTimeNoSeconds = "2006-01-02 15:04"
)

// NowUnix returns the current Unix timestamp in seconds.
func NowUnix() int64 {
	return time.Now().Unix()
}

// FromUnix converts a seconds-based Unix timestamp to time.Time.
func FromUnix(unix int64) time.Time {
	return time.Unix(unix, 0)
}

// NowTimestamp returns the current timestamp in milliseconds.
func NowTimestamp() int64 {
	return Timestamp(time.Now())
}

// Timestamp returns the milliseconds timestamp for t.
func Timestamp(t time.Time) int64 {
	return t.UnixNano() / 1e6
}

// FromTimestamp converts a milliseconds timestamp to time.Time.
func FromTimestamp(timestamp int64) time.Time {
	return time.Unix(0, timestamp*int64(time.Millisecond))
}

// Format formats a time according to the provided layout.
func Format(time time.Time, layout string) string {
	return time.Format(layout)
}

// Parse parses a time string according to layout.
func Parse(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// GetDay return yyyyMMdd
func GetDay(time time.Time) int {
	ret, _ := strconv.Atoi(time.Format("20060102"))
	return ret
}

// WithTimeAsStartOfDay returns the start of the day for the given time.
func WithTimeAsStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func WithTimeAsEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 1000000000-1, t.Location())
}
