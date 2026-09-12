package domain

import (
	"testing"
	"time"
)

// The cases that were wrong while the day was measured in UTC: for a UTC-3
// student the UTC day flips at 21:00 local, right in the middle of homework
// time.
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
			// The instant is what matters, not the zone it is expressed in.
			name: "an instant stored as UTC still lands on the local day",
			from: time.Date(2026, 8, 11, 1, 0, 0, 0, time.UTC),  // 22:00 of the 10th local
			to:   time.Date(2026, 8, 11, 23, 0, 0, 0, time.UTC), // 20:00 of the 11th local
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
