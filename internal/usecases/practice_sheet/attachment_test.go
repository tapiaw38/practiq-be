package practicesheet

import "testing"

// scoreWith mirrors the accounting the submit loop does, so the rule "a file
// nobody graded must not count as wrong" is pinned down.
func scoreWith(results []attachmentOutcome) (correct, total int, score float64, allPending bool) {
	total = len(results)
	for _, outcome := range results {
		switch {
		case outcome.NeedsReview:
			total--
		case outcome.IsCorrect:
			correct++
		}
	}
	if total > 0 {
		score = float64(correct) / float64(total) * 100
	}
	allPending = total <= 0 && len(results) > 0
	return
}

func TestPendingAttachmentsAreExcludedFromScore(t *testing.T) {
	t.Run("a pending file does not drag the score down", func(t *testing.T) {
		// One correct answer plus an ungradeable PDF should read as 100%,
		// not 50%.
		_, total, score, allPending := scoreWith([]attachmentOutcome{
			{IsCorrect: true},
			{NeedsReview: true},
		})
		if total != 1 || score != 100 {
			t.Errorf("expected 1 graded answer at 100%%, got total=%d score=%v", total, score)
		}
		if allPending {
			t.Error("some answers were graded, so the sheet is not fully pending")
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

	t.Run("everything pending yields no score to act on", func(t *testing.T) {
		_, total, score, allPending := scoreWith([]attachmentOutcome{
			{NeedsReview: true},
			{NeedsReview: true},
		})
		if total != 0 || score != 0 {
			t.Errorf("expected nothing graded, got total=%d score=%v", total, score)
		}
		if !allPending {
			t.Error("a sheet with only pending answers must be flagged, or the student fails at 0%")
		}
	})

	t.Run("an empty submission is not pending review", func(t *testing.T) {
		if _, _, _, allPending := scoreWith(nil); allPending {
			t.Error("no answers at all is not the same as awaiting review")
		}
	})
}
