package grade

import (
	"reflect"
	"testing"
)

func TestFilterAllowedUserIDs(t *testing.T) {
	tests := map[string]struct {
		requesterID string
		isTeacher   bool
		userIDs     []string
		want        []string
	}{
		"teacher sees every requested id": {
			requesterID: "teacher-1",
			isTeacher:   true,
			userIDs:     []string{"student-1", "student-2", "teacher-1"},
			want:        []string{"student-1", "student-2", "teacher-1"},
		},
		"student sees only themselves out of the batch": {
			requesterID: "student-1",
			isTeacher:   false,
			userIDs:     []string{"student-1", "student-2"},
			want:        []string{"student-1"},
		},
		"student not in the batch gets nothing, not an error": {
			requesterID: "student-3",
			isTeacher:   false,
			userIDs:     []string{"student-1", "student-2"},
			want:        nil,
		},
		"empty batch stays empty for a teacher": {
			requesterID: "teacher-1",
			isTeacher:   true,
			userIDs:     []string{},
			want:        []string{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := filterAllowedUserIDs(tc.requesterID, tc.isTeacher, tc.userIDs)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("filterAllowedUserIDs() = %#v, want %#v", got, tc.want)
			}
		})
	}
}
