package practicesheet

import "testing"

// What happens to an answer the assistant was asked to judge and could not.
//
// Before this, the string comparison left standing decided on its own, so a
// practice written "cuatro" against an expected "4" came back wrong — scored,
// final, and with nobody to appeal to, because a practice has no teacher queue.
func TestUnresolvedEvaluation(t *testing.T) {
	t.Run("an exact match stands as correct", func(t *testing.T) {
		// Nothing to interpret: the words are the words. Sending this to a
		// teacher, or dropping it from the score, would punish the student for
		// an outage that changed nothing about their answer.
		isCorrect, ungraded := unresolvedEvaluation(true)
		if !isCorrect || ungraded {
			t.Fatalf("isCorrect=%v ungraded=%v, want true/false", isCorrect, ungraded)
		}
	})

	t.Run("a mismatch is indeterminate, not wrong", func(t *testing.T) {
		// The comparison cannot tell a different wording from a different
		// answer. Only the assistant could, and it did not reply.
		isCorrect, ungraded := unresolvedEvaluation(false)
		if isCorrect {
			t.Fatal("a mismatch the assistant never judged must not be scored as correct")
		}
		if !ungraded {
			t.Fatal("a mismatch nobody judged must be left ungraded, not marked wrong")
		}
	})
}

// The message has to match what actually happens next, or it sends the student
// after a correction nobody is going to make.
func TestEvaluationUnavailableFeedbackMatchesWhoGrades(t *testing.T) {
	if got := evaluationUnavailableFeedback(true); got == evaluationUnavailableFeedback(false) {
		t.Fatal("a level test is reviewed by a teacher and a practice is not; the two must not say the same thing")
	}
}
