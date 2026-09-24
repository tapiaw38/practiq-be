package practicesheet

import (
	"testing"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

// The gate used to reject every scheduled sheet: blocked before the date and
// "expired" from the date onwards, so a scheduled level test could never be
// taken. These cases pin the window down.
func TestSheetWindowState(t *testing.T) {
	now := time.Now()

	hourAgo := time.Now().Add(-time.Hour)
	hourAhead := time.Now().Add(time.Hour)
	twoHoursAgo := time.Now().Add(-2 * time.Hour)

	cases := []struct {
		name  string
		sheet domain.PracticeSheet
		want  windowState
	}{
		{
			name:  "no schedule is always open",
			sheet: domain.PracticeSheet{},
			want:  windowOpen,
		},
		{
			name:  "before the opening date",
			sheet: domain.PracticeSheet{ScheduledAt: &hourAhead},
			want:  windowNotYetOpen,
		},
		{
			// The regression: this used to report closed, so a scheduled test
			// was unreachable from the moment it opened.
			name:  "open once the date passes and no closing date",
			sheet: domain.PracticeSheet{ScheduledAt: &hourAgo},
			want:  windowOpen,
		},
		{
			name:  "inside the window",
			sheet: domain.PracticeSheet{ScheduledAt: &hourAgo, AvailableUntil: &hourAhead},
			want:  windowOpen,
		},
		{
			name:  "after the window closed",
			sheet: domain.PracticeSheet{ScheduledAt: &twoHoursAgo, AvailableUntil: &hourAgo},
			want:  windowClosed,
		},
		{
			name:  "the closing instant is closed",
			sheet: domain.PracticeSheet{ScheduledAt: &twoHoursAgo, AvailableUntil: &now},
			want:  windowClosed,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := sheetWindowState(&tc.sheet, now); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
