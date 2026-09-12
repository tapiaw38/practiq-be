package domain

import (
	"strings"
	"testing"
)

func TestStudentsToDeactivate(t *testing.T) {
	// Least recently active first, which is the order the caller supplies.
	byActivity := []string{"dormant", "old", "recent", "active", "today"}

	cases := []struct {
		name   string
		max    int
		chosen []string
		want   string
	}{
		{
			name: "a plan with room deactivates nobody",
			max:  5,
			want: "",
		},
		{
			name: "a bigger plan deactivates nobody",
			max:  9,
			want: "",
		},
		{
			// The default: the ones who have not practised in longest go
			// first, and whoever never practised goes before them.
			name: "without a choice the least recently active go",
			max:  2,
			want: "dormant,old,recent",
		},
		{
			// The teacher's call wins, even when it keeps someone dormant.
			name:   "a choice is honoured",
			max:    2,
			chosen: []string{"dormant", "today"},
			want:   "old,recent,active",
		},
		{
			// Choosing fewer than the plan allows leaves room, and the
			// automatic order fills it rather than cutting more than needed.
			name:   "a partial choice is filled by activity",
			max:    3,
			chosen: []string{"today"},
			want:   "dormant,old",
		},
		{
			// Asking to keep more than the plan allows cannot be honoured as
			// given; the extra falls back to the order.
			name:   "a choice bigger than the plan is truncated",
			max:    2,
			chosen: []string{"dormant", "old", "recent", "active"},
			want:   "recent,active,today",
		},
		{
			name: "a plan of zero deactivates everyone",
			max:  0,
			want: "dormant,old,recent,active,today",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := strings.Join(StudentsToDeactivate(byActivity, tc.max, tc.chosen), ",")
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

// Whatever is deactivated, what is left must fit the plan. A rule that cuts too
// few leaves a school over its limit; one that cuts too many takes students
// away for nothing.
func TestWhatSurvivesAlwaysFitsThePlan(t *testing.T) {
	byActivity := []string{"a", "b", "c", "d", "e", "f"}

	for max := 0; max <= len(byActivity)+2; max++ {
		for _, chosen := range [][]string{nil, {"f"}, {"a", "f"}, {"a", "b", "c", "d", "e", "f"}} {
			out := StudentsToDeactivate(byActivity, max, chosen)
			surviving := len(byActivity) - len(out)
			if surviving > max {
				t.Fatalf("max=%d chosen=%v left %d students, over the plan", max, chosen, surviving)
			}
			if max <= len(byActivity) && surviving != max {
				t.Fatalf("max=%d chosen=%v left %d students, want %d", max, chosen, surviving, max)
			}
		}
	}
}
