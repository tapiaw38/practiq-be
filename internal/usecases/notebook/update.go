package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, error)
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

func (u *updateUsecase) Execute(ctx context.Context, id string, input UpdateInput) (*UpdateOutput, error) {
	app := u.contextFactory()
	if err := app.Repositories.Notebook.Update(ctx, id, domain.Notebook{
		Title:       input.Title,
		Description: input.Description,
	}); err != nil {
		return nil, err
	}
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UpdateOutput{Data: toNotebookData(nb)}, nil
}
