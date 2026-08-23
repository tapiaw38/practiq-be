package notebook

import "testing"

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

// The teacher is called in only where the assistant fell short: it never read
// the answer, or it read it and could not decide.
func TestSubmissionNeedsTeacherReview(t *testing.T) {
	correct := true
	incorrect := false

	for _, test := range []struct {
		name           string
		hasStudentWork bool
		aiIsCorrect    *bool
		want           bool
	}{
		{name: "ocr failed", hasStudentWork: true, want: true},
		{name: "unreadable", hasStudentWork: true, want: true},
		{name: "evaluation failed", hasStudentWork: true, want: true},
		{name: "no assistant configured", hasStudentWork: true, want: true},
		{name: "graded correct", hasStudentWork: true, aiIsCorrect: &correct},
		{name: "graded incorrect", hasStudentWork: true, aiIsCorrect: &incorrect},
		{name: "nothing submitted"},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := submissionNeedsTeacherReview(test.hasStudentWork, test.aiIsCorrect); got != test.want {
				t.Fatalf("submissionNeedsTeacherReview() = %t, want %t", got, test.want)
			}
		})
	}
}

// An unverified transcription of the teacher's statement used to send every
// submission on that page to the teacher, even the ones graded cleanly.
func TestUnverifiedStatementDoesNotForceReview(t *testing.T) {
	graded := true
	if submissionNeedsTeacherReview(true, &graded) {
		t.Fatal("a clean verdict stands on its own, whatever the statement's state")
	}
}
