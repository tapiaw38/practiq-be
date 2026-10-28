package attemptreview

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestToReviewDataDoesNotExposeCanonicalStatementMediaURL(t *testing.T) {
	data := toReviewData(domain.PendingAttemptReview{
		StatementMediaURL: "https://private-bucket.example/exercises/teacher/prompt.mp3",
	})

	if data.StatementMediaViewURL != "" {
		t.Fatalf("canonical statement media URL must only be exposed after presigning, got %q", data.StatementMediaViewURL)
	}
}
