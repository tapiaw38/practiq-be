package practicesheet

import (
	"testing"
	"time"
)

func TestParseScheduledAt(t *testing.T) {
	t.Run("empty clears the schedule", func(t *testing.T) {
		at, err := parseScheduledAt("  ")
		if err != nil || at != nil {
			t.Errorf("expected no schedule, got at=%v err=%v", at, err)
		}
	})

	t.Run("offset is normalized to UTC", func(t *testing.T) {
		// The browser sends the teacher's local time with its offset; storing it
		// as-is would make the gate fire at the wrong moment.
		at, err := parseScheduledAt("2026-09-01T14:30:00-03:00")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := time.Date(2026, 9, 1, 17, 30, 0, 0, time.UTC)
		if !at.Equal(want) {
			t.Errorf("expected %v, got %v", want, at.UTC())
		}
		if at.Location() != time.UTC {
			t.Errorf("expected UTC, got %v", at.Location())
		}
	})

	t.Run("invalid values are rejected", func(t *testing.T) {
		for _, value := range []string{"2026-09-01 14:30", "mañana", "2026-13-01T00:00:00Z"} {
			if _, err := parseScheduledAt(value); err == nil {
				t.Errorf("expected an error for %q", value)
			}
		}
	})
}
