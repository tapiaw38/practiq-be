package domain

import (
	"strings"
	"testing"
)

func TestStudentsToDeactivate(t *testing.T) {

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

			name: "without a choice the least recently active go",
			max:  2,
			want: "dormant,old,recent",
		},
		{

			name:   "a choice is honoured",
			max:    2,
			chosen: []string{"dormant", "today"},
			want:   "old,recent,active",
		},
		{

			name:   "a partial choice is filled by activity",
			max:    3,
			chosen: []string{"today"},
			want:   "dormant,old",
		},
		{

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
