package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, courseID string) (*ListOutput, error)
	}

	ListOutput struct {
		Data []NotebookData `json:"data"`
	}

	listUsecase struct{ contextFactory appcontext.Factory }
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, courseID string) (*ListOutput, error) {
	app := u.contextFactory()
	notebooks, err := app.Repositories.Notebook.List(ctx, courseID)
	if err != nil {
		return nil, err
	}
	data := make([]NotebookData, 0, len(notebooks))
	for _, nb := range notebooks {
		nb := nb
		data = append(data, toNotebookData(&nb))
	}
	return &ListOutput{Data: data}, nil
}
