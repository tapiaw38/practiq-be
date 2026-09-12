package practicesheet

import "github.com/tapiaw38/practiq-be/internal/domain"

const exerciseTypeFillBlanks = "fill_blanks"

// blanksAnswersMatch grades a fill-in-the-blanks answer by exact comparison.
// The rules live in the domain so exercise creation validates against the same
// definition the grader uses.
func blanksAnswersMatch(student, expected string) bool {
	return domain.BlanksAnswersMatch(student, expected)
}
