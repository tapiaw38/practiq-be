package domain

import "testing"

func TestExerciseMediaURL(t *testing.T) {
	cases := []struct {
		name     string
		metadata string
		want     string
	}{
		{"no metadata", "", ""},
		{"no media key", `{"accept":["audio"]}`, ""},
		{"reads the url", `{"accept":["audio"],"media_url":"https://s3/x.mp4"}`, "https://s3/x.mp4"},
		{"broken json is not media", `{oops`, ""},
		{"wrong type is not media", `{"media_url":42}`, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Exercise{Metadata: tc.metadata}.MediaURL()
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}
