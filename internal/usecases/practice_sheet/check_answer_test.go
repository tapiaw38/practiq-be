package practicesheet

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestCheckableExercise(t *testing.T) {
	cases := []struct {
		name     string
		exercise domain.Exercise
		want     bool
	}{
		{"plain text", domain.Exercise{Type: "text"}, true},
		{"fill blanks", domain.Exercise{Type: exerciseTypeFillBlanks}, true},
		{"unparseable metadata is not media", domain.Exercise{Type: "text", Metadata: "{"}, true},
		{"attachment", domain.Exercise{Type: exerciseTypeAttachment}, false},
		{
			"statement carries media",
			domain.Exercise{Type: "text", Metadata: `{"media_url":"https://cdn/audio.mp3"}`},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := checkableExercise(tc.exercise); got != tc.want {
				t.Fatalf("checkableExercise = %v, want %v", got, tc.want)
			}
		})
	}
}
