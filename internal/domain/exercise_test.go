package domain

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestMetadataWithoutTeacherImageKeepsEverythingElse(t *testing.T) {
	const drawing = "data:image/png;base64,iVBORw0KGgo="

	cases := []struct {
		name     string
		metadata string
		wantKeys []string
	}{
		{
			name:     "the drawing goes, the configuration stays",
			metadata: `{"teacher_image":"` + drawing + `","media_url":"https://b/x.png","options":["a"]}`,
			wantKeys: []string{"media_url", "options"},
		},
		{
			// The key was renamed twice; a payload saved under an old name has
			// to be stripped too or the image ships anyway.
			name:     "older key names are stripped as well",
			metadata: `{"imageData":"` + drawing + `","layout":"grid"}`,
			wantKeys: []string{"layout"},
		},
		{
			name:     "metadata without a drawing is untouched",
			metadata: `{"options":["a","b"]}`,
			wantKeys: []string{"options"},
		},
		{
			name:     "empty metadata stays empty",
			metadata: "",
			wantKeys: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Exercise{Metadata: tc.metadata}.MetadataWithoutTeacherImage()
			if strings.Contains(got, "base64") {
				t.Fatalf("the drawing survived: %s", got)
			}
			if tc.metadata == "" {
				if got != "" {
					t.Fatalf("got %q, want empty", got)
				}
				return
			}
			var values map[string]json.RawMessage
			if err := json.Unmarshal([]byte(got), &values); err != nil {
				t.Fatalf("result is not valid JSON: %v", err)
			}
			if len(values) != len(tc.wantKeys) {
				t.Fatalf("got keys %v, want %v", values, tc.wantKeys)
			}
			for _, key := range tc.wantKeys {
				if _, ok := values[key]; !ok {
					t.Fatalf("lost key %q", key)
				}
			}
		})
	}
}

// A bare data URL was a valid metadata value before the field became JSON.
// Returning it unchanged would defeat the whole point of stripping.
func TestMetadataWithoutTeacherImageFailsClosedOnUnparseableValues(t *testing.T) {
	got := Exercise{Metadata: "data:image/png;base64,iVBORw0KGgo="}.MetadataWithoutTeacherImage()
	if got != "{}" {
		t.Fatalf("got %q, want {}", got)
	}
}

func TestTeacherImageReadsBothStorageForms(t *testing.T) {
	stored := "https://bucket.s3.amazonaws.com/exercises/a.png"
	if got := (Exercise{Metadata: `{"teacher_image":"` + stored + `"}`}).TeacherImage(); got != stored {
		t.Fatalf("bucket URL: got %q, want %q", got, stored)
	}

	inline := "data:image/png;base64,iVBORw0KGgo="
	if got := (Exercise{Metadata: `{"teacher_image":"` + inline + `"}`}).TeacherImage(); got != inline {
		t.Fatalf("data URL: got %q, want %q", got, inline)
	}

	// A value that is neither is not an image, and handing it to a fetcher would
	// be following an arbitrary string out of the database.
	if got := (Exercise{Metadata: `{"teacher_image":"/etc/passwd"}`}).TeacherImage(); got != "" {
		t.Fatalf("got %q, want empty", got)
	}
}
