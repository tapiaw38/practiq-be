package attemptreview

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ReviewUsecase interface {
		Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError)
	}

	reviewUsecase struct {
		contextFactory appcontext.Factory
	}

	ReviewInput struct {
		IsCorrect *bool  `json:"is_correct" binding:"required"`
		Feedback  string `json:"feedback"`
	}

	ReviewOutput struct {
		Data OperationResultData `json:"data"`
	}
)

func NewReviewUsecase(contextFactory appcontext.Factory) ReviewUsecase {
	return &reviewUsecase{contextFactory: contextFactory}
}

func (u *reviewUsecase) Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool, input ReviewInput) (*ReviewOutput, apperrors.ApplicationError) {
	if input.IsCorrect == nil {
		return nil, apperrors.NewBadRequestError("is_correct is required")
	}
	isCorrect := *input.IsCorrect
	app := u.contextFactory()

	owner, err := app.Repositories.StudentAttempt.GetTeacherForAttempt(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if owner == "" {
		return nil, apperrors.NewNotFoundError("attempt not found")
	}
	if !isSuperAdmin && owner != teacherID {
		return nil, apperrors.NewForbiddenError()
	}

	attemptCtx, err := app.Repositories.StudentAttempt.GetAttemptContext(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if attemptCtx.SheetType != sheetTypeLevelTest {
		return nil, apperrors.NewBadRequestError("only level test answers are corrected by the teacher")
	}

	if err := app.Repositories.StudentAttempt.Review(ctx, attemptID, isCorrect, input.Feedback); err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}

	applyLevelTestOutcome(ctx, app, attemptID)

	return &ReviewOutput{Data: OperationResultData{Message: "attempt reviewed"}}, nil
}
