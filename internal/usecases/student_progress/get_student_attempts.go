package studentprogress

import (
	"context"
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type AttemptData struct {
	ID              string `json:"id"`
	StudentID       string `json:"student_id"`
	ExerciseID      string `json:"exercise_id"`
	PracticeSheetID string `json:"practice_sheet_id"`
	AnswerText      string `json:"answer_text"`
	IsCorrect       bool   `json:"is_correct"`
	Score           float64 `json:"score"`
	TimeSpentSecs   int    `json:"time_spent_seconds"`
	HintsUsed       int    `json:"hints_used"`
	CreatedAt       string `json:"created_at"`
}

type AttemptListOutput struct {
	Data []AttemptData `json:"data"`
}

func toAttemptData(a domain.StudentAttempt) AttemptData {
	return AttemptData{
		ID:              a.ID,
		StudentID:       a.StudentID,
		ExerciseID:      a.ExerciseID,
		PracticeSheetID: a.PracticeSheetID,
		AnswerText:      a.AnswerText,
		IsCorrect:       a.IsCorrect,
		Score:           a.Score,
		TimeSpentSecs:   a.TimeSpentSecs,
		HintsUsed:       a.HintsUsed,
		CreatedAt:       a.CreatedAt.Format(time.RFC3339),
	}
}

type GetStudentAttemptsUsecase interface {
	Execute(ctx context.Context, studentID, sheetID string) (*AttemptListOutput, apperrors.ApplicationError)
}

type getStudentAttemptsUsecase struct {
	factory appcontext.Factory
}

func NewGetStudentAttemptsUsecase(factory appcontext.Factory) GetStudentAttemptsUsecase {
	return &getStudentAttemptsUsecase{factory: factory}
}

func (u *getStudentAttemptsUsecase) Execute(ctx context.Context, studentID, sheetID string) (*AttemptListOutput, apperrors.ApplicationError) {
	app := u.factory()

	attempts, err := app.Repositories.StudentAttempt.ListBySheet(ctx, studentID, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []AttemptData
	for _, a := range attempts {
		data = append(data, toAttemptData(a))
	}
	if data == nil {
		data = []AttemptData{}
	}

	return &AttemptListOutput{Data: data}, nil
}
