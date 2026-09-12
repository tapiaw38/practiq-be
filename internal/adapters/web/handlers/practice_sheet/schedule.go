package practicesheet

import (
	"errors"
	"strings"
	"time"
)

// parseScheduledAt accepts an RFC 3339 timestamp and normalizes it to UTC. An
// empty value clears the schedule.
func parseScheduledAt(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errors.New("scheduled_at must be an RFC 3339 timestamp, e.g. 2026-09-01T14:30:00Z")
	}
	utc := parsed.UTC()
	return &utc, nil
}
