package ai

import (
	"strings"
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestBuildHelpPromptReviewAnswerIsDirect(t *testing.T) {
	prompt := buildHelpPrompt("review_answer", "Revisá mi respuesta", "10", &domain.Exercise{Question: "2 + 3 + 5", CorrectAnswer: "10"}, "", nil)
	for _, expected := range []string{"Modo REVISIÓN", "Respuesta declarada por estudiante: 10", "Correcta."} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q: %s", expected, prompt)
		}
	}
}

func TestBuildHelpPromptFillBlanksUsesPuzzleLanguage(t *testing.T) {
	exercise := &domain.Exercise{
		Type:          "fill_blanks",
		Question:      "1 + {{1}} = 6",
		CorrectAnswer: `{"1":"5"}`,
		Metadata:      `{"options":["5","4","3"]}`,
	}
	prompt := buildHelpPrompt("review_answer", "Revisá mi respuesta actual.", `{"1":"4"}`, exercise, "1er Grado", nil)
	for _, expected := range []string{"rompecabezas", "[hueco 1]", "Bloques visibles: 5, 4, 3", "Hueco 1: 4", "Hueco 1: 5"} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("prompt missing %q: %s", expected, prompt)
		}
	}
	if strings.Contains(prompt, `{"1":"`) {
		t.Fatalf("prompt must not expose fill blanks JSON: %s", prompt)
	}
}

func TestBuildHelpPromptHintIsBounded(t *testing.T) {
	prompt := buildHelpPrompt("hint", "Dame una pista", "", &domain.Exercise{Question: "2 + 3 + 5"}, "", nil)
	if !strings.Contains(prompt, "Máximo 20 palabras") {
		t.Fatalf("hint prompt must require short output: %s", prompt)
	}
}
