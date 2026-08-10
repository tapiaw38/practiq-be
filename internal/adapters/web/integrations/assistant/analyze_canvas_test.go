package assistant

import "testing"

func TestCanvasResponseRejectsTutorVerdicts(t *testing.T) {
	for _, response := range []string{"Correcta.", "incorrecto", "¡Incorrecta!", "Es correcta."} {
		if isExpectedCanvasResponse(response) {
			t.Fatalf("expected tutor verdict %q to be rejected", response)
		}
	}
}

func TestCanvasResponseAcceptsStudentAnswerOrUnreadable(t *testing.T) {
	for _, response := range []string{"4", "10", "x = 5", "UNREADABLE"} {
		if !isExpectedCanvasResponse(response) {
			t.Fatalf("expected canvas response %q to be accepted", response)
		}
	}
}
