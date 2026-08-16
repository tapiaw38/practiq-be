package practicesheet

import (
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// Measured in UTC, the day flipped at 21:00 in Argentina: practising twice one
// evening counted as two days, and practising two evenings running counted as
// one. Both are the common case, since homework happens after dinner.
func TestCalcStreakUsesTheStudentsDay(t *testing.T) {
	loc := domain.StudentLocation(domain.DefaultTimezone)
	now := time.Now().In(loc)
	at := func(daysAgo, hour int) *time.Time {
		d := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, loc).
			AddDate(0, 0, -daysAgo)
		return &d
	}

	cases := []struct {
		name    string
		current *domain.StudentTopicProgress
		want    int
	}{
		{
			name:    "the first practice starts the streak",
			current: nil,
			want:    1,
		},
		{
			name:    "practising again the same day does not add one",
			current: &domain.StudentTopicProgress{StreakDays: 5, LastPracticedAt: at(0, 2)},
			want:    5,
		},
		{
			name:    "practising the next day adds one",
			current: &domain.StudentTopicProgress{StreakDays: 5, LastPracticedAt: at(1, 22)},
			want:    6,
		},
		{
			name:    "a gap starts over",
			current: &domain.StudentTopicProgress{StreakDays: 5, LastPracticedAt: at(4, 9)},
			want:    1,
		},
		{
			name:    "a progress row that was never practised starts the streak",
			current: &domain.StudentTopicProgress{StreakDays: 3},
			want:    1,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := calcStreak(tc.current, loc); got != tc.want {
				t.Fatalf("got %d, want %d", got, tc.want)
			}
		})
	}
}
