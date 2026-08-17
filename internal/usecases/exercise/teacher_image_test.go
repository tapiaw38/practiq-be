package exercise

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

type fakeStorage struct {
	storage.NoopImageStorage
	uploaded string
	err      error
}

func (f *fakeStorage) UploadDataURI(_ context.Context, folder, _, dataURI string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.uploaded = dataURI
	return "https://bucket.s3.amazonaws.com/" + folder + "/stored.png", nil
}

const inlineDrawing = "data:image/png;base64,iVBORw0KGgo="

func TestStoreTeacherImage(t *testing.T) {
	const storedURL = "https://bucket.s3.amazonaws.com/exercises/old.png"

	cases := []struct {
		name     string
		incoming string
		previous domain.Exercise
		want     string
		wantUp   bool
	}{
		{
			name:     "an inline drawing is uploaded and replaced by its URL",
			incoming: `{"teacher_image":"` + inlineDrawing + `"}`,
			want:     "https://bucket.s3.amazonaws.com/exercises/stored.png",
			wantUp:   true,
		},
		{
			// Responses no longer carry the drawing, so this is what the editor
			// sends after opening an exercise and saving it unchanged. Losing
			// the statement here would be the whole feature breaking quietly.
			name:     "a save that omits the drawing keeps the stored one",
			incoming: `{"options":["a"]}`,
			previous: domain.Exercise{Metadata: `{"teacher_image":"` + storedURL + `"}`},
			want:     storedURL,
		},
		{
			name:     "an explicit empty value clears the drawing",
			incoming: `{"teacher_image":""}`,
			previous: domain.Exercise{Metadata: `{"teacher_image":"` + storedURL + `"}`},
			want:     "",
		},
		{
			name:     "a value already stored is not uploaded again",
			incoming: `{"teacher_image":"` + storedURL + `"}`,
			previous: domain.Exercise{Metadata: `{"teacher_image":"` + storedURL + `"}`},
			want:     storedURL,
		},
		{
			name:     "an exercise that never had one stays without it",
			incoming: `{"options":["a"]}`,
			want:     "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store := &fakeStorage{}
			app := &appcontext.Context{ImageStorage: store}

			got := storeTeacherImage(context.Background(), app, "teacher-1", tc.incoming, tc.previous)

			if image := (domain.Exercise{Metadata: got}).TeacherImage(); image != tc.want {
				t.Fatalf("stored image = %q, want %q", image, tc.want)
			}
			if strings.Contains(got, "base64") {
				t.Fatalf("the drawing stayed inline: %s", got)
			}
			if uploaded := store.uploaded != ""; uploaded != tc.wantUp {
				t.Fatalf("uploaded = %v, want %v", uploaded, tc.wantUp)
			}
		})
	}
}

// A storage failure must not throw away what the teacher just drew.
func TestStoreTeacherImageKeepsTheDrawingWhenUploadFails(t *testing.T) {
	app := &appcontext.Context{ImageStorage: &fakeStorage{err: errors.New("s3 is down")}}

	got := storeTeacherImage(context.Background(), app, "teacher-1",
		`{"teacher_image":"`+inlineDrawing+`"}`, domain.Exercise{})

	if image := (domain.Exercise{Metadata: got}).TeacherImage(); image != inlineDrawing {
		t.Fatalf("got %q, want the drawing kept inline", image)
	}
}
