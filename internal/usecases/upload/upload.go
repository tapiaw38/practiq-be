package upload

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/storage"
)

type (
	Usecase interface {
		Execute(context.Context, Input) (*Output, apperrors.ApplicationError)
	}

	usecase struct {
		contextFactory appcontext.Factory
	}

	Input struct {
		UserID string
		// Folder groups objects in the bucket, e.g. "attachments".
		Folder      string
		Filename    string
		ContentType string
		Reader      io.Reader
		Size        int64
	}

	Output struct {
		Data FileData `json:"data"`
	}

	FileData struct {
		URL         string `json:"url"`
		Filename    string `json:"filename"`
		ContentType string `json:"content_type"`
		Kind        string `json:"kind"`
		Size        int64  `json:"size"`
	}
)

func NewUsecase(contextFactory appcontext.Factory) Usecase {
	return &usecase{contextFactory: contextFactory}
}

func (u *usecase) Execute(ctx context.Context, input Input) (*Output, apperrors.ApplicationError) {
	app := u.contextFactory()

	if app.ImageStorage == nil || !app.ImageStorage.IsConfigured() {
		return nil, apperrors.NewApplicationError(mappings.UploadNotConfiguredError, nil)
	}
	if input.Size > storage.MaxUploadBytes {
		return nil, apperrors.NewApplicationError(mappings.UploadTooLargeError, nil)
	}

	// Read one byte past the cap so a lying Content-Length is still caught.
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

	// Resolve first so the response reports the type the file is actually
	// stored with, which may be the sniffed one.
	contentType, kind, _, err := storage.ResolveContentType(input.ContentType, body)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
	}

	url, err := app.ImageStorage.UploadFile(ctx, folder, input.UserID, input.Filename, contentType, body)
	if err != nil {
		if errors.Is(err, storage.ErrUnsupportedFileType) {
			return nil, apperrors.NewApplicationError(mappings.UploadUnsupportedTypeError, err)
		}
		log.Printf("[upload] failed user_id=%s filename=%q err=%v", input.UserID, input.Filename, err)
		return nil, apperrors.NewApplicationError(mappings.UploadError, err)
	}

	return &Output{Data: FileData{
		URL:         url,
		Filename:    input.Filename,
		ContentType: contentType,
		Kind:        string(kind),
		Size:        int64(len(body)),
	}}, nil
}
