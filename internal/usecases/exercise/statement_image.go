package exercise

import (
	"context"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	// StatementImageUsecase serves the handwritten statement of one exercise.
	//
	// It exists so no listing or sheet payload has to carry the drawing. It also
	// hides where the drawing lives: exercises saved by the current editor point
	// at the bucket, older ones keep the image inline, and both come back here
	// as plain bytes.
	//
	// Serving the bytes through the API rather than handing out a bucket URL is
	// what lets a browser draw the image into a canvas. A cross-origin image
	// taints the canvas, and both callers that matter — the editor loading a
	// drawing to change it, and the practice screen composing the page it sends
	// for grading — export that canvas afterwards.
	StatementImageUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, exerciseID string) (*StatementImageOutput, apperrors.ApplicationError)
	}

	statementImageUsecase struct {
		contextFactory appcontext.Factory
	}

	StatementImageOutput struct {
		Content     []byte
		ContentType string
	}
)

var errMalformedDataURL = errors.New("statement image is not a base64 data URL")

func NewStatementImageUsecase(contextFactory appcontext.Factory) StatementImageUsecase {
	return &statementImageUsecase{contextFactory: contextFactory}
}

func (u *statementImageUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, exerciseID string) (*StatementImageOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	ex, err := app.Repositories.Exercise.Get(ctx, exerciseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	if ex == nil {
		return nil, apperrors.NewNotFoundError("exercise not found")
	}

	// Seeing the statement is seeing the course it belongs to: the same rule the
	// exercise listing applies, so this opens nothing the requester could not
	// already read.
	if appErr := requesterCanReadExercise(ctx, app, requesterID, isAdmin, ex.TopicID); appErr != nil {
		return nil, appErr
	}

	image := ex.TeacherImage()
	if image == "" {
		return nil, apperrors.NewNotFoundError("exercise has no statement image")
	}

	if strings.HasPrefix(image, "data:") {
		content, contentType, decodeErr := decodeDataURL(image)
		if decodeErr != nil {
			return nil, apperrors.NewApplicationError(mappings.ExerciseListError, decodeErr)
		}
		return &StatementImageOutput{Content: content, ContentType: contentType}, nil
	}

	if app.ImageStorage == nil {
		return nil, apperrors.NewApplicationError(mappings.UploadNotConfiguredError, nil)
	}
	content, contentType, err := app.ImageStorage.FetchFile(ctx, image)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ExerciseListError, err)
	}
	return &StatementImageOutput{Content: content, ContentType: contentType}, nil
}

// decodeDataURL unpacks the inline form older exercises were saved in.
func decodeDataURL(value string) ([]byte, string, error) {
	comma := strings.Index(value, ",")
	if comma < 0 {
		return nil, "", errMalformedDataURL
	}
	header := value[len("data:"):comma]
	contentType, isBase64 := strings.CutSuffix(header, ";base64")
	if !isBase64 {
		return nil, "", errMalformedDataURL
	}
	content, err := base64.StdEncoding.DecodeString(value[comma+1:])
	if err != nil {
		return nil, "", err
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	return content, contentType, nil
}
