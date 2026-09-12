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
	// StatementImageUsecase serves the handwritten statement of the exercise an
	// attempt belongs to.
	//
	// It is a separate call on purpose: the image is a base64 canvas of a few
	// hundred KB, and inlining one per row would turn the review queue into a
	// multi-megabyte response for a list the teacher only skims.
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
		// Image is a data URL, empty when the exercise has no handwritten
		// statement.
		Image string `json:"image"`
	}
)

func NewStatementImageUsecase(contextFactory appcontext.Factory) StatementImageUsecase {
	return &statementImageUsecase{contextFactory: contextFactory}
}

func (u *statementImageUsecase) Execute(ctx context.Context, attemptID, teacherID string, isSuperAdmin bool) (*StatementImageOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Same ownership rule as reviewing: whoever may grade the attempt may see
	// the statement it was answering.
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

// asDataURL keeps this endpoint's contract while the storage underneath changed.
//
// Exercises saved by the current editor keep the drawing in the bucket, so the
// value here is a URL the reviewer's browser cannot open on its own. Older ones
// already hold a data URL and pass straight through.
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
