package domain

import (
	"testing"
	"time"
)

func TestDaysBetweenUsesTheStudentsDay(t *testing.T) {
	loc := StudentLocation(DefaultTimezone)

	cases := []struct {
		name     string
		from, to time.Time
		want     int
	}{
		{
			name: "two sessions the same evening are the same day",
			from: time.Date(2026, 8, 10, 20, 0, 0, 0, loc),
			to:   time.Date(2026, 8, 10, 22, 0, 0, 0, loc),
			want: 0,
		},
		{
			name: "two consecutive evenings are one day apart",
			from: time.Date(2026, 8, 10, 22, 0, 0, 0, loc),
			to:   time.Date(2026, 8, 11, 20, 0, 0, 0, loc),
			want: 1,
		},
		{
			name: "two consecutive mornings are one day apart",
			from: time.Date(2026, 8, 10, 9, 0, 0, 0, loc),
			to:   time.Date(2026, 8, 11, 9, 0, 0, 0, loc),
			want: 1,
		},
		{
			name: "a gap is more than one day",
			from: time.Date(2026, 8, 10, 9, 0, 0, 0, loc),
			to:   time.Date(2026, 8, 14, 9, 0, 0, 0, loc),
			want: 4,
		},
		{

			name: "an instant stored as UTC still lands on the local day",
			from: time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC),
			to:   time.Date(2026, 8, 11, 23, 0, 0, 0, time.UTC),
			want: 1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := DaysBetween(tc.from, tc.to, loc); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}

func TestStudentLocationFallsBack(t *testing.T) {
	if StudentLocation("").String() != DefaultTimezone {
		t.Fatalf("empty timezone should fall back to the product default")
	}
	if StudentLocation("Not/AZone").String() != DefaultTimezone {
		t.Fatalf("an unknown timezone should fall back rather than fail")
	}
	if StudentLocation("America/Mexico_City").String() != "America/Mexico_City" {
		t.Fatalf("a valid timezone should be honoured")
	}
}
