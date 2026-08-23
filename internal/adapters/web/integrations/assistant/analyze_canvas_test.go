package assistant

import (
	"strings"
	"testing"
)

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

// A page of worked problems is the normal case here and is far past the 120
// characters the single-answer check allows, which is what made every notebook
// transcription look malformed and retry until the attempts ran out.
func TestNotebookResponseAcceptsAFullPage(t *testing.T) {
	page := "1) 2 + 3 = 5\n2) 4 + 4 = 8\n3) 10 - 7 = 3\n4) 6 + 6 = 12\n" +
		"5) 9 - 2 = 7\n6) 8 + 5 = 13\n7) 12 - 4 = 8\n8) 3 + 9 = 12\n9) 15 - 6 = 9"
	if len(page) <= 120 {
		t.Fatalf("fixture must exceed the single-answer limit, got %d chars", len(page))
	}
	if isExpectedCanvasResponse(page) {
		t.Fatal("the single-answer check should still reject a full page")
	}
	if !isExpectedNotebookResponse(page) {
		t.Fatal("expected a full notebook page to be accepted")
	}
}

func TestNotebookResponseKeepsTheOtherGuards(t *testing.T) {
	if !isExpectedNotebookResponse("UNREADABLE") {
		t.Fatal("UNREADABLE must stay acceptable")
	}
	for _, response := range []string{
		"",
		"Correcta.",
		"incorrecto",
		`{"answer": "4"}`,
		strings.Repeat("a", maxNotebookTranscriptionChars+1),
	} {
		if isExpectedNotebookResponse(response) {
			t.Fatalf("expected notebook response %q to be rejected", truncateForLog(response, 40))
		}
	}
}

// normalizeCanvasResponse unwraps fences before any check sees them, so a
// fenced transcription is usable rather than malformed. Pinned because the
// original validator carried a Contains("```") guard that could never fire and
// read as protection that was not there.
func TestNotebookResponseUnwrapsFencedTranscription(t *testing.T) {
	if !isExpectedNotebookResponse("```\n1) 2 + 3 = 5\n```") {
		t.Fatal("a fenced transcription should be unwrapped and accepted")
	}
	if got := normalizeCanvasResponse("```\n1) 2 + 3 = 5\n```"); strings.Contains(got, "`") {
		t.Fatalf("fences should be stripped, got %q", got)
	}
}

// The teacher's page image is uploaded and stored as a URL, so the expected
// answer reached the model as "https://….png" instead of a placeholder.
func TestNotebookPromptCarriesContextNotAURL(t *testing.T) {
	prompt := buildNotebookCanvasPrompt("Cuaderno - Pagina 1. Titulo: Sumas. Instrucciones: resolve.")
	if !strings.Contains(prompt, "Titulo: Sumas") {
		t.Fatal("prompt must carry the page context")
	}
	if !strings.Contains(prompt, unreadableResponse) {
		t.Fatal("prompt must keep the UNREADABLE escape hatch")
	}
	if strings.Contains(prompt, "respuesta correcta esperada") {
		t.Fatal("the notebook prompt must not ask for a single expected answer")
	}
	if empty := buildNotebookCanvasPrompt("   "); !strings.Contains(empty, "sin contexto") {
		t.Fatal("a blank context must degrade to a placeholder")
	}
}

func TestUnreadableIsAcceptedButRetryable(t *testing.T) {
	if !isExpectedNotebookResponse(unreadableResponse) {
		t.Fatal("UNREADABLE must remain a well-formed response")
	}
	if !isExpectedCanvasResponse(unreadableResponse) {
		t.Fatal("UNREADABLE must remain well-formed for the single-answer check too")
	}
	if !strings.EqualFold(normalizeCanvasResponse("  unreadable  "), unreadableResponse) {
		t.Fatal("the retry branch compares against the normalized value, so it must fold case and spacing")
	}
}
