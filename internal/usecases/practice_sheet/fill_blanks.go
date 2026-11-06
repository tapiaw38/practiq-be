package practicesheet

import "github.com/tapiaw38/practiq-be/internal/domain"

const exerciseTypeFillBlanks = "fill_blanks"

func blanksAnswersMatch(student, expected string) bool {
	return domain.BlanksAnswersMatch(student, expected)
}
