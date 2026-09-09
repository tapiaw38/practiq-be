package domain

import "testing"

func TestPersonalSchoolName(t *testing.T) {
	cases := []struct {
		name        string
		teacherName string
		want        string
	}{
		{
			name:        "the teacher's name makes the school's",
			teacherName: "Ana García",
			want:        "Escuela de Ana García",
		},
		{
			// identity.FullName falls back to the user id when auth-api-be has
			// no name, and an id is not a name to show anyone.
			name:        "no name keeps the placeholder",
			teacherName: "",
			want:        PlaceholderSchoolName,
		},
		{
			// Half a name — "Escuela de " with nothing after — reads as a bug
			// on a screen, so whitespace counts as no name.
			name:        "whitespace is no name",
			teacherName: "   ",
			want:        PlaceholderSchoolName,
		},
		{
			name:        "surrounding spaces are trimmed",
			teacherName: "  Ana García  ",
			want:        "Escuela de Ana García",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := PersonalSchoolName(tc.teacherName); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
