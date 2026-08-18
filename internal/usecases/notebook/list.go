package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*ListOutput, apperrors.ApplicationError)
	}

	ListOutput struct {
		Data []NotebookData `json:"data"`
	}

	listUsecase struct{ contextFactory appcontext.Factory }
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if appErr := requesterCanReadCourse(ctx, app, requesterID, isSuperAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	notebooks, err := app.Repositories.Notebook.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	data := make([]NotebookData, 0, len(notebooks))
	for _, nb := range notebooks {
		nb := nb
		resolveNotebookImages(ctx, app, &nb)
		data = append(data, toNotebookData(&nb))
	}
	return &ListOutput{Data: data}, nil
}
