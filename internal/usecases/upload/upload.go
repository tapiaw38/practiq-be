package upload

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

const previewLinkTTL = time.Hour

type (
	UploadUsecase interface {
		Execute(context.Context, UploadInput) (*UploadOutput, apperrors.ApplicationError)
	}

	uploadUsecase struct {
		contextFactory appcontext.Factory
	}

	UploadInput struct {
		UserID string

		Folder      string
		Filename    string
		ContentType string
		Reader      io.Reader
		Size        int64
	}

	UploadOutput struct {
		Data FileData `json:"data"`
	}

	FileData struct {
		URL string `json:"url"`

		PreviewURL  string `json:"preview_url,omitempty"`
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Kind        string `json:"kind"`
		Size        int64  `json:"size"`
	}
)

func NewUploadUsecase(contextFactory appcontext.Factory) UploadUsecase {
	return &uploadUsecase{contextFactory: contextFactory}
}

func (u *uploadUsecase) Execute(ctx context.Context, input UploadInput) (*UploadOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if app.ImageStorage == nil || !app.ImageStorage.IsConfigured() {
		return nil, apperrors.NewApplicationError(mappings.UploadNotConfiguredError, nil)
	}
	if input.Size > storage.MaxUploadBytes {
		return nil, apperrors.NewApplicationError(mappings.UploadTooLargeError, nil)
	}

	body, err := io.ReadAll(io.LimitReader(input.Reader, storage.MaxUploadBytes+1))
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UploadError, err)
	}
	if len(body) > storage.MaxUploadBytes {
		return nil, apperrors.NewApplicationError(mappings.UploadTooLargeError, nil)
	}

	folder := input.Folder
	if folder == "" {
		folder = "uploads"
	}

	contentType, kind, _, err := storage.ResolveContentType(input.ContentType, body)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
	}
	if folder == "exercises" && !storage.AllowedAsStatementMaterial(kind) {
		return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, storage.ErrUnsupportedFileType)
	}

	url, err := app.ImageStorage.UploadFile(ctx, folder, input.UserID, input.Filename, contentType, body)
	if err != nil {
		if errors.Is(err, storage.ErrUnsupportedFileType) {
			return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
		}
		log.Printf("[upload] failed user_id=%s filename=%q err=%v", input.UserID, input.Filename, err)
		return nil, apperrors.NewApplicationError(mappings.UploadError, err)
	}

	previewURL, _ := app.ImageStorage.PresignGetURL(url, previewLinkTTL)

	return &UploadOutput{Data: FileData{
		URL:         url,
		PreviewURL:  previewURL,
		Filename:    input.Filename,
		ContentType: contentType,
		Kind:        string(kind),
		Size:        int64(len(body)),
	}}, nil
}
