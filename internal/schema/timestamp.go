package schema

import (
	"encoding/json"
	"time"
)

const (
	dquote          = '"'
	timestampLayout = "2026-06-12T15:20:00.123Z"
)

var (
	_ json.Marshaler   = (*Timestamp)(nil)
	_ json.Unmarshaler = (*Timestamp)(nil)
)

// Timestamp is a wrapper around [time.Time] that marshals and unmarshals from JSON as an ISO 8601 string. By Outline convention, the time is always in UTC (with the Z suffix).
type Timestamp struct {
	inner time.Time
}

// Creates a new timestamp.
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{t.UTC()}
}

// Creates a new timestamp with the current time.
func NewTimestampNow() Timestamp {
	return NewTimestamp(time.Now())
}

// Returns the time value.
func (t *Timestamp) Time() time.Time {
	return t.inner
}

// Marshals the timestamp to JSON as a simplifed ISO 8601 string.
//
// This format is equivalent to the output of Javascript's `Date.prototype.toISOString()`.
// See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString.
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timestampLayout)+2)
	b = append(b, dquote)
	b = t.inner.AppendFormat(b, timestampLayout)
	b = append(b, dquote)
	return b, nil
}

// Unmarshals the timestamp from JSON.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	return t.inner.UnmarshalJSON(data)
}
