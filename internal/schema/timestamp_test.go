package schema

import (
	"encoding/json"
	"testing"
	"time"
)

func TestTimestampMarshal(t *testing.T) {
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "zero value",
			time:     time.Time{},
			expected: `"0001-01-01T00:00:00.000Z"`,
		},
		{
			name: "this year",
			// I'm feeling festive 123 miliseconds into Christmas!
			time:     time.Date(2026, time.December, 25, 0, 0, 0, 123000000, time.UTC),
			expected: `"2026-12-25T00:00:00.123Z"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := NewTimestamp(tt.time)
			b, err := json.Marshal(&ts)
			if err != nil {
				t.Errorf("failed to marshal timestamp: %s", err)
			}

			s := string(b)
			if s != tt.expected {
				t.Errorf("expected timestamp %q, got %q", tt.expected, s)
			}
		})
	}
}
