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

func TestReviewFlagSurvivesEveryFailurePath(t *testing.T) {
	page := &domain.NotebookPage{
		ContentData:   "https://bucket.s3.amazonaws.com/image/notebook/x/y.png",
		StatementText: "3+2=",
	}

	start := statementNeedsTeacherReview(page)
	if !start {
		t.Fatal("an unverified image statement must start out flagged")
	}

	for _, path := range []struct {
		name        string
		failed      bool
		hasAnswer   bool
		wantsReview bool
	}{
		{"ocr failed", true, false, true},
		{"unreadable", true, false, true},
		{"evaluation failed", true, true, true},
		{"evaluated fine", false, true, true},
		{"nothing submitted", false, false, false},
	} {
		needsReview := start
		if path.failed {
			needsReview = true
		}
		if !path.hasAnswer && !path.failed {
			needsReview = false
		}
		if needsReview != path.wantsReview {
			t.Fatalf("%s: needsReview=%v, want %v", path.name, needsReview, path.wantsReview)
		}
	}
}

func TestVerifiedStatementNeverFlagsOnSuccess(t *testing.T) {
	page := &domain.NotebookPage{
		ContentData:       "https://bucket.s3.amazonaws.com/image/notebook/x/y.png",
		StatementText:     "3+2=",
		StatementVerified: true,
	}
	if statementNeedsTeacherReview(page) {
		t.Fatal("a teacher-verified statement must not force review on a clean run")
	}
}
