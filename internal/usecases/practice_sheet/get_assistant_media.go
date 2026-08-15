package practicesheet

import (
	"context"
	"fmt"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

type (
	GetAssistantMediaUsecase interface {
		Execute(context.Context, string, bool, string, string) (*AssistantMediaOutput, apperrors.ApplicationError)
	}

	getAssistantMediaUsecase struct{ contextFactory appcontext.Factory }

	AssistantMediaOutput struct {
		Content     []byte
		ContentType string
	}
)

func NewGetAssistantMediaUsecase(contextFactory appcontext.Factory) GetAssistantMediaUsecase {
	return &getAssistantMediaUsecase{contextFactory: contextFactory}
}

// Execute serves only media belonging to a sheet the requester can open. This
// avoids browser-to-bucket CORS requirements and never accepts an arbitrary URL.
func (u *getAssistantMediaUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, sheetID, exerciseID string) (*AssistantMediaOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	ps, err := app.Repositories.PracticeSheet.Get(ctx, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}
	if appErr := requesterCanReadCourse(ctx, app, requesterID, isAdmin, ps.CourseID); appErr != nil {
		return nil, appErr
	}
	if appErr := ensureSheetIsOpen(ctx, app, ps, requesterID, isAdmin); appErr != nil {
		return nil, appErr
	}

	var mediaURL string
	for _, item := range ps.Exercises {
		if item.Exercise.ID == exerciseID {
			mediaURL = item.Exercise.MediaURL()
			break
		}
	}
	if mediaURL == "" {
		return nil, apperrors.NewNotFoundError("exercise media not found")
	}
	if app.ImageStorage == nil {
		return nil, apperrors.NewApplicationError(mappings.UploadNotConfiguredError, nil)
	}
	body, contentType, err := app.ImageStorage.FetchFile(ctx, mediaURL)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	kind, _, err := storage.ClassifyContentType(contentType)
	if err != nil || (kind != storage.FileKindImage && kind != storage.FileKindAudio) {
		return nil, apperrors.NewBadRequestError(fmt.Sprintf("exercise media is not supported by the assistant: %s", contentType))
	}
	return &AssistantMediaOutput{Content: body, ContentType: contentType}, nil
}
