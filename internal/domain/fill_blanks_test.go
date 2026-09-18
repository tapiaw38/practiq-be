package domain

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestBlankIDsInStatement(t *testing.T) {
	t.Run("returns the ids in the order they appear", func(t *testing.T) {
		ids, err := BlankIDsInStatement("En {{2}} se firmó en {{1}}")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(ids) != 2 || ids[0] != 2 || ids[1] != 1 {
			t.Fatalf("got %v", ids)
		}
	})

	t.Run("rejects a repeated marker", func(t *testing.T) {
		_, err := BlankIDsInStatement("{{1}} y otra vez {{1}}")
		if !errors.Is(err, ErrDuplicateMarker) {
			t.Fatalf("expected ErrDuplicateMarker, got %v", err)
		}
	})

	t.Run("rejects a statement without blanks", func(t *testing.T) {
		if _, err := BlankIDsInStatement("Sin huecos"); !errors.Is(err, ErrNoBlanks) {
			t.Fatalf("expected ErrNoBlanks, got %v", err)
		}
	})
}

func TestParseBlanksAnswerRejectsCollidingKeys(t *testing.T) {
	// "01" and "1" name the same blank. Deciding between them would depend on
	// Go's randomised map order, so the same answer would grade differently on
	// each attempt. Run it enough times that a coin flip could not pass.
	for i := 0; i < 100; i++ {
		if _, err := ParseBlanksAnswer(`{"01":"a","1":"b"}`); !errors.Is(err, ErrDuplicateAnswerID) {
			t.Fatalf("expected ErrDuplicateAnswerID, got %v", err)
		}
	}

	if !BlanksAnswersMatch(`{"01":"a"}`, `{"1":"a"}`) {
		t.Fatal("a padded key that collides with nothing should still match")
	}
	if BlanksAnswersMatch(`{"01":"a","1":"b"}`, `{"1":"a"}`) {
		t.Fatal("a colliding answer must never be graded as correct")
	}
	for i := 0; i < 100; i++ {
		if BlanksAnswersMatch(`{"01":"","1":"a"}`, `{"1":"a"}`) {
			t.Fatal("an empty colliding answer must never be graded as correct")
		}
	}
}

func metadataFor(t *testing.T, blanks map[int]string, options []string) string {
	t.Helper()
	type blank struct {
		ID     int    `json:"id"`
		Answer string `json:"answer"`
	}
	list := make([]blank, 0, len(blanks))
	for id, answer := range blanks {
		list = append(list, blank{ID: id, Answer: answer})
	}
	encoded, err := json.Marshal(map[string]any{"blanks": list, "options": options, "layout": "text"})
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return string(encoded)
}

func TestValidateFillBlanksExercise(t *testing.T) {
	t.Run("accepts a solvable exercise", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "1816", 2: "Tucumán"}, []string{"1816", "Tucumán", "1810"})
		if err := ValidateFillBlanksExercise("En {{1}} se firmó en {{2}}", metadata); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("accepts two blanks sharing an answer when the pool has both blocks", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "verbo", 2: "verbo"}, []string{"verbo", "verbo", "sujeto"})
		if err := ValidateFillBlanksExercise("{{1}} y {{2}}", metadata); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects a repeated answer with a single block", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "verbo", 2: "verbo"}, []string{"verbo", "sujeto"})
		if err := ValidateFillBlanksExercise("{{1}} y {{2}}", metadata); !errors.Is(err, ErrMissingOptions) {
			t.Fatalf("expected ErrMissingOptions, got %v", err)
		}
	})

	t.Run("rejects a blank with no answer", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "  "}, []string{"algo"})
		if err := ValidateFillBlanksExercise("Completá {{1}}", metadata); !errors.Is(err, ErrEmptyBlankAnswer) {
			t.Fatalf("expected ErrEmptyBlankAnswer, got %v", err)
		}
	})

	t.Run("rejects blanks that do not match the statement", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "algo", 3: "otra"}, []string{"algo", "otra"})
		if err := ValidateFillBlanksExercise("{{1}} y {{2}}", metadata); !errors.Is(err, ErrBlanksMismatch) {
			t.Fatalf("expected ErrBlanksMismatch, got %v", err)
		}
	})

	t.Run("rejects a statement with a repeated marker", func(t *testing.T) {
		metadata := metadataFor(t, map[int]string{1: "algo"}, []string{"algo"})
		if err := ValidateFillBlanksExercise("{{1}} y {{1}}", metadata); !errors.Is(err, ErrDuplicateMarker) {
			t.Fatalf("expected ErrDuplicateMarker, got %v", err)
		}
	})

	t.Run("rejects unusable metadata", func(t *testing.T) {
		if err := ValidateFillBlanksExercise("Completá {{1}}", "no soy json"); err == nil {
			t.Fatal("expected an error")
		}
		if err := ValidateFillBlanksExercise("Completá {{1}}", ""); err == nil {
			t.Fatal("expected an error")
		}
	})
}
