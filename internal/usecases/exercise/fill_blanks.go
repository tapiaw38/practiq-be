package exercise

import (
	"github.com/tapiaw38/practiq-be/internal/domain"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
)

const typeFillBlanks = "fill_blanks"

func validateFillBlanks(exerciseType, question, metadata, correctAnswer string) apperrors.ApplicationError {
	if exerciseType != typeFillBlanks {
		return nil
	}
	if err := domain.ValidateFillBlanksExercise(question, metadata); err != nil {
		return apperrors.NewBadRequestError(err.Error())
	}
	if err := domain.ValidateFillBlanksGradingAnswer(correctAnswer, metadata); err != nil {
		return apperrors.NewBadRequestError(err.Error())
	}
	return nil
}
