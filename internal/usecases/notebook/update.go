package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	UpdateInput struct {
		Title       string
		Description string
	}

	UpdateOutput struct {
		Data NotebookData `json:"data"`
	}

	updateUsecase struct{ contextFactory appcontext.Factory }
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify notebook exists and check ownership
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if nb == nil {
		return nil, apperrors.NewNotFoundError("notebook not found")
	}

	if !isAdmin && nb.TeacherID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	if err := app.Repositories.Notebook.Update(ctx, id, domain.Notebook{
		Title:       input.Title,
		Description: input.Description,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}

	nb, err = app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}

	return &UpdateOutput{Data: toNotebookData(nb)}, nil
}
