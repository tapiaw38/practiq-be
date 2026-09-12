package practicesheet

import (
	"testing"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

func TestValidateLevelTestAttemptsRequiresExactExerciseSet(t *testing.T) {
	exercises := []domain.PracticeSheetExercise{
		{Exercise: domain.Exercise{ID: "first"}},
		{Exercise: domain.Exercise{ID: "second"}},
	}

	cases := []struct {
		name     string
		attempts []AttemptInput
		wantErr  bool
	}{
		{"complete set", []AttemptInput{{ExerciseID: "first"}, {ExerciseID: "second"}}, false},
		{"partial set", []AttemptInput{{ExerciseID: "first"}}, true},
		{"duplicate replaces an exercise", []AttemptInput{{ExerciseID: "first"}, {ExerciseID: "first"}}, true},
		{"foreign exercise", []AttemptInput{{ExerciseID: "first"}, {ExerciseID: "other"}}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateLevelTestAttempts(exercises, tc.attempts)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}
