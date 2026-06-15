package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, id, studentID string) (*GetOutput, error)
	}

	GetOutput struct {
		Data NotebookData `json:"data"`
	}

	getUsecase struct{ contextFactory appcontext.Factory }
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, id, studentID string) (*GetOutput, error) {
	app := u.contextFactory()
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if nb == nil {
		return nil, nil
	}

	// Attach student submissions to each page
	if studentID != "" {
		for i := range nb.Pages {
			sub, _ := app.Repositories.Notebook.GetSubmission(ctx, nb.Pages[i].ID, studentID)
			nb.Pages[i].Submission = sub
		}
	}

	return &GetOutput{Data: toNotebookData(nb)}, nil
}
