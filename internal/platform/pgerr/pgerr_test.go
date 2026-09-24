package pgerr

import (
	"errors"
	"fmt"
	"testing"

	"github.com/lib/pq"
)

func TestIsUniqueViolation(t *testing.T) {
	dup := &pq.Error{Code: "23505", Constraint: "idx_grades_school_name"}

	cases := []struct {
		name       string
		err        error
		constraint string
		want       bool
	}{
		{
			name:       "the duplicate it was asked about",
			err:        dup,
			constraint: "idx_grades_school_name",
			want:       true,
		},
		{
			// The reason the constraint is named: a table carries more than
			// one unique index, and "that name is taken" would be the wrong
			// answer for a different one.
			name:       "a duplicate on another index",
			err:        dup,
			constraint: "idx_subjects_school_name",
		},
		{
			name:       "a different postgres error",
			err:        &pq.Error{Code: "23503", Constraint: "idx_grades_school_name"},
			constraint: "idx_grades_school_name",
		},
		{
			name:       "still found once wrapped",
			err:        fmt.Errorf("updating grade: %w", dup),
			constraint: "idx_grades_school_name",
			want:       true,
		},
		{
			name:       "a plain error",
			err:        errors.New("connection reset"),
			constraint: "idx_grades_school_name",
		},
		{
			name:       "no error at all",
			err:        nil,
			constraint: "idx_grades_school_name",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsUniqueViolation(tc.err, tc.constraint); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
