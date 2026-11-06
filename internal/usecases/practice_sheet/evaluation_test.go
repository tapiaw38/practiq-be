package practicesheet

import "testing"

func TestUnresolvedEvaluation(t *testing.T) {
	t.Run("an exact match stands as correct", func(t *testing.T) {

		isCorrect, ungraded := unresolvedEvaluation(true)
		if !isCorrect || ungraded {
			t.Fatalf("isCorrect=%v ungraded=%v, want true/false", isCorrect, ungraded)
		}
	})

	t.Run("a mismatch is indeterminate, not wrong", func(t *testing.T) {

		isCorrect, ungraded := unresolvedEvaluation(false)
		if isCorrect {
			t.Fatal("a mismatch the assistant never judged must not be scored as correct")
		}
		if !ungraded {
			t.Fatal("a mismatch nobody judged must be left ungraded, not marked wrong")
		}
	})
}

func TestEvaluationUnavailableFeedbackMatchesWhoGrades(t *testing.T) {
	if got := evaluationUnavailableFeedback(true); got == evaluationUnavailableFeedback(false) {
		t.Fatal("a level test is reviewed by a teacher and a practice is not; the two must not say the same thing")
	}
}
