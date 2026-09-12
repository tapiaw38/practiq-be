package practicesheet

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestNeedsReviewForStatementMedia(t *testing.T) {
	withMedia := domain.Exercise{Metadata: `{"media_url":"https://storage.example/exercises/teacher/prompt.mp3"}`}
	withoutMedia := domain.Exercise{}

	for _, test := range []struct {
		name                              string
		exercise                          domain.Exercise
		hasText, hasCanvas, hasAttachment bool
		want                              bool
	}{
		{name: "unanswered media exercise", exercise: withMedia, want: false},
		{name: "text answer with media", exercise: withMedia, hasText: true, want: true},
		{name: "canvas answer with media", exercise: withMedia, hasCanvas: true, want: true},
		{name: "attachment answer with media", exercise: withMedia, hasAttachment: true, want: true},
		{name: "text answer without media", exercise: withoutMedia, hasText: true, want: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := needsReviewForStatementMedia(test.exercise, test.hasText, test.hasCanvas, test.hasAttachment); got != test.want {
				t.Fatalf("needsReviewForStatementMedia() = %t, want %t", got, test.want)
			}
		})
	}
}
