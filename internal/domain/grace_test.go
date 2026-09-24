package domain

import (
	"testing"
	"time"
)

func TestInGrace(t *testing.T) {
	paidUntil := time.Date(2026, 10, 24, 2, 49, 15, 0, time.UTC)

	cases := []struct {
		name string
		now  time.Time
		want bool
	}{
		{"still paid is not grace", paidUntil.Add(-time.Hour), false},
		{"the moment it lapses", paidUntil, true},
		{"four days later", paidUntil.AddDate(0, 0, 4), true},
		{"a minute before the fifth day is up", GraceEndsAt(paidUntil).Add(-time.Minute), true},
		{"the fifth day is up", GraceEndsAt(paidUntil), false},
		{"long gone", paidUntil.AddDate(0, 1, 0), false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := InGrace(paidUntil, tc.now); got != tc.want {
				t.Fatalf("InGrace = %v, want %v", got, tc.want)
			}
		})
	}
}

// A subscription that never recorded a paid period cannot be in grace: there
// is no day to count five from, and treating it as grace would hand the plan
// to anybody whose dates failed to sync.
func TestNoPaidPeriodIsNotGrace(t *testing.T) {
	if InGrace(time.Time{}, time.Now()) {
		t.Fatal("a zero paid-until must never grant grace")
	}
}
