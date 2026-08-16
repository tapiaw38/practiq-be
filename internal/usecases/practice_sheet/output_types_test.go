package practicesheet

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestStudentMetadataHidesSolutions(t *testing.T) {
	const full = `{"media_url":"https://s3/x.png","layout":"text",` +
		`"blanks":[{"id":1,"answer":"sol"},{"id":2,"answer":"luna"}],` +
		`"options":["sol","luna","mar"]}`

	got := studentMetadata(full)

	if strings.Contains(got, "media_url") {
		t.Fatalf("the canonical bucket URL leaked: %s", got)
	}
	// The secret is the blank-to-option mapping, not the words themselves:
	// the option pool legitimately contains every answer, shuffled, because the
	// student picks from it. What must not ship is which one goes where.
	if strings.Contains(got, `"answer"`) {
		t.Fatalf("a blank answer leaked: %s", got)
	}

	var values map[string]any
	if err := json.Unmarshal([]byte(got), &values); err != nil {
		t.Fatalf("result is not JSON: %v", err)
	}
	// What the student still needs to render the exercise.
	if values["layout"] != "text" {
		t.Fatalf("layout was dropped: %s", got)
	}
	options, ok := values["options"].([]any)
	if !ok || len(options) != 3 {
		t.Fatalf("the option pool was dropped: %s", got)
	}
	blanks, ok := values["blanks"].([]any)
	if !ok || len(blanks) != 2 {
		t.Fatalf("blank ids were dropped: %s", got)
	}
	first, ok := blanks[0].(map[string]any)
	if !ok || first["id"] != float64(1) {
		t.Fatalf("blank id was dropped: %s", got)
	}
}

func TestStudentMetadataEdgeCases(t *testing.T) {
	cases := []struct {
		name     string
		metadata string
		want     string
	}{
		{"empty stays empty", "", ""},
		{"broken json fails closed", `{oops`, `{}`},
		{"no blanks key", `{"accept":["audio"]}`, `{"accept":["audio"]}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := studentMetadata(tc.metadata); got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}

	t.Run("an unexpected blanks shape is dropped, not forwarded", func(t *testing.T) {
		got := studentMetadata(`{"blanks":"sol,luna"}`)
		if strings.Contains(got, "sol") {
			t.Fatalf("unparsable blanks were forwarded: %s", got)
		}
	})

	t.Run("malformed metadata cannot leak an answer", func(t *testing.T) {
		got := studentMetadata(`{"blanks":[{"id":1,"answer":"secreto"}]`)
		if strings.Contains(got, "secreto") {
			t.Fatalf("malformed metadata leaked an answer: %s", got)
		}
	})
}
