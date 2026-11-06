package attemptreview

import (
	"context"
	"encoding/base64"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	StatementImageUsecase interface {
		Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool) (*StatementImageOutput, apperrors.ApplicationError)
	}

	statementImageUsecase struct {
		contextFactory appcontext.Factory
	}

	StatementImageOutput struct {
		Data StatementImageData `json:"data"`
	}

	StatementImageData struct {
		Image string `json:"image"`
	}
)

func NewStatementImageUsecase(contextFactory appcontext.Factory) StatementImageUsecase {
	return &statementImageUsecase{contextFactory: contextFactory}
}

func (u *statementImageUsecase) Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool) (*StatementImageOutput, apperrors.ApplicationError) {
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

	exerciseID, err := app.Repositories.StudentAttempt.GetExerciseIDForAttempt(ctx, attemptID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AttemptReviewError, err)
	}
	if exerciseID == "" {
		return nil, apperrors.NewNotFoundError("attempt not found")
	}

	exercise, err := app.Repositories.Exercise.Get(ctx, exerciseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if exercise == nil {
		return nil, apperrors.NewNotFoundError("exercise not found")
	}

	return &StatementImageOutput{
		Data: StatementImageData{Image: u.asDataURL(ctx, app, exercise.TeacherImage())},
	}, nil
}

func (u *statementImageUsecase) asDataURL(ctx context.Context, app *appcontext.Context, image string) string {
	if image == "" || strings.HasPrefix(image, "data:") {
		return image
	}
	if app.ImageStorage == nil {
		return ""
	}
	content, contentType, err := app.ImageStorage.FetchFile(ctx, image)
	if err != nil || len(content) == 0 {
		return ""
	}
	if contentType == "" {
		contentType = "image/png"
	}
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(content)
}
