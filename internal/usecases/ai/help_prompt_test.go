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

func TestBuildHelpPromptHintIsBounded(t *testing.T) {
	prompt := buildHelpPrompt("hint", "Dame una pista", "", &domain.Exercise{Question: "2 + 3 + 5"}, "", nil)
	if !strings.Contains(prompt, "Máximo 20 palabras") {
		t.Fatalf("hint prompt must require short output: %s", prompt)
	}
}
