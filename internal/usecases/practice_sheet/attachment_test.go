package practicesheet

import (
	"strings"
	"testing"
)

// scoreWith mirrors the accounting the submit loop does, so the rule "an answer
// nobody graded must not count as wrong" is pinned down.
func scoreWith(results []attachmentOutcome) (correct, total int, score float64, allUngraded bool) {
	total = len(results)
	for _, outcome := range results {
		switch {
		case outcome.Ungraded:
			total--
		case outcome.IsCorrect:
			correct++
		}
	}
	if total > 0 {
		score = float64(correct) / float64(total) * 100
	}
	allUngraded = total <= 0 && len(results) > 0
	return
}

// A graded file answer must not block the student: it scores right away and
// records who graded it, so a teacher can review it later without the practice
// flow waiting on them.
func TestAssistantGradedAnswerScoresImmediately(t *testing.T) {
	approved := true
	outcome := attachmentOutcome{
		IsCorrect:          true,
		Feedback:           "Muy buena lectura",
		Ungraded:           false,
		AISuggestedCorrect: &approved,
	}

	if outcome.Ungraded {
		t.Error("a graded answer must not hold the student up")
	}
	if outcome.AISuggestedCorrect == nil {
		t.Error("the assistant's verdict must be recorded so the teacher can review it")
	}

	correct, total, score, allUngraded := scoreWith([]attachmentOutcome{outcome})
	if correct != 1 || total != 1 || score != 100 {
		t.Errorf("expected it to count as correct, got %d/%d at %v", correct, total, score)
	}
	if allUngraded {
		t.Error("a graded answer is not ungraded")
	}
}

func TestUngradedAttachmentsAreExcludedFromScore(t *testing.T) {
	t.Run("an ungraded file does not drag the score down", func(t *testing.T) {
		// One correct answer plus an ungradeable PDF should read as 100%,
		// not 50%.
		_, total, score, allUngraded := scoreWith([]attachmentOutcome{
			{IsCorrect: true},
			{Ungraded: true},
		})
		if total != 1 || score != 100 {
			t.Errorf("expected 1 graded answer at 100%%, got total=%d score=%v", total, score)
		}
		if allUngraded {
			t.Error("some answers were graded, so the sheet is not fully ungraded")
		}
	})

	t.Run("a wrong answer still counts", func(t *testing.T) {
		correct, total, score, _ := scoreWith([]attachmentOutcome{
			{IsCorrect: true},
			{IsCorrect: false},
		})
		if correct != 1 || total != 2 || score != 50 {
			t.Errorf("expected 1/2 at 50%%, got %d/%d at %v", correct, total, score)
		}
	})

	t.Run("nothing graded yields no score to act on", func(t *testing.T) {
		_, total, score, allUngraded := scoreWith([]attachmentOutcome{
			{Ungraded: true},
			{Ungraded: true},
		})
		if total != 0 || score != 0 {
			t.Errorf("expected nothing graded, got total=%d score=%v", total, score)
		}
		if !allUngraded {
			t.Error("a sheet with no graded answers must be flagged, or the student fails at 0%")
		}
	})

	t.Run("an empty submission is not ungraded", func(t *testing.T) {
		if _, _, _, allUngraded := scoreWith(nil); allUngraded {
			t.Error("no answers at all is not the same as an answer nobody could grade")
		}
	})
}

// The teacher only corrects level tests: the homework notebook has its own
// queue, and a practice must resolve on submit instead of waiting on anyone.
func TestOnlyLevelTestsGoToTheTeacher(t *testing.T) {
	if !teacherGradesSheet(sheetTypeLevelTest) {
		t.Error("a level test decides promotion, so a teacher confirms it")
	}
	if teacherGradesSheet("practice") {
		t.Error("a practice must never be queued for the teacher")
	}
	if teacherGradesSheet("") {
		t.Error("an unknown sheet type must not reach the teacher's queue")
	}
}

// An answer nobody graded is explained differently on each sheet: only a level
// test can promise a correction.
func TestUngradedFeedbackOnlyPromisesReviewOnALevelTest(t *testing.T) {
	if got := ungradedAttachmentFeedback(true); got == ungradedAttachmentFeedback(false) {
		t.Error("a practice must not be told a teacher will correct it")
	}
	if got := statementMediaFeedback(true); got == statementMediaFeedback(false) {
		t.Error("a practice must not be told a teacher will correct it")
	}
	// "pedí revisión" pointed the student at a correction a practice no longer
	// offers.
	if strings.Contains(unreadableCanvasFeedback(false), "revisión") ||
		strings.Contains(unreadableCanvasFeedback(false), "docente") {
		t.Error("a practice must not send the student after a review nobody will do")
	}
	if !strings.Contains(unreadableCanvasFeedback(true), "docente") {
		t.Error("on a level test the unreadable answer does go to the teacher")
	}
}

// A handwritten answer nobody transcribed used to fall through to the string
// comparison and score as wrong against an empty answer. On a level test that
// cost the student the promotion for a transcription that never ran.
func TestUntranscribedHandwritingIsNotWrong(t *testing.T) {
	for _, test := range []struct {
		name                                     string
		canvasAwaitsOCR, hasText, assistantReady bool
		want                                     bool
	}{
		{name: "no assistant configured", canvasAwaitsOCR: true, want: true},
		{name: "assistant read it", canvasAwaitsOCR: true, hasText: true, assistantReady: true},
		{name: "assistant available, still blank", canvasAwaitsOCR: true, assistantReady: true},
		{name: "typed answer needs no transcription", hasText: true},
		{name: "statement media keeps OCR out of it", hasText: true, assistantReady: true},
		{name: "nothing was answered"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := transcriptionUnavailable(test.canvasAwaitsOCR, test.hasText, test.assistantReady)
			if got != test.want {
				t.Fatalf("transcriptionUnavailable() = %t, want %t", got, test.want)
			}
		})
	}
}
