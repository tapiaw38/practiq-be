package notebook

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestStatementNeedsTeacherReview(t *testing.T) {
	imageURL := "https://bucket.s3.amazonaws.com/image/notebook/x/y.png"

	if !statementNeedsTeacherReview(&domain.NotebookPage{ContentData: imageURL}) {
		t.Fatal("an unverified transcription must keep the verdict a suggestion")
	}
	if statementNeedsTeacherReview(&domain.NotebookPage{ContentData: imageURL, StatementVerified: true}) {
		t.Fatal("once the teacher verified the statement the verdict stands on its own")
	}
	if statementNeedsTeacherReview(&domain.NotebookPage{ContentData: "Resolve las sumas"}) {
		t.Fatal("a text statement was never guessed at, so it needs no review")
	}
	if statementNeedsTeacherReview(nil) {
		t.Fatal("a missing page must not flag review")
	}
}

func TestPageHasImageStatement(t *testing.T) {
	for _, value := range []string{
		"https://bucket.s3.amazonaws.com/a/b.png",
		"data:image/png;base64,iVBORw0KGgo",
	} {
		if !pageHasImageStatement(value) {
			t.Fatalf("expected %q to count as an image statement", value)
		}
	}
	for _, value := range []string{"", "   ", "Resolve las sumas", "https://example.com/page"} {
		if pageHasImageStatement(value) {
			t.Fatalf("expected %q not to count as an image statement", value)
		}
	}
}
