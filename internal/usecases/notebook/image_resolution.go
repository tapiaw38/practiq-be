package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/usecases/datauri"
)

func resolveNotebookImages(ctx context.Context, app *appcontext.Context, nb *domain.Notebook) {
	if nb == nil {
		return
	}
	for i := range nb.Pages {
		nb.Pages[i].ContentData = datauri.Resolve(ctx, app, nb.Pages[i].ContentData)
		if nb.Pages[i].Submission != nil {
			nb.Pages[i].Submission.CanvasData = datauri.Resolve(ctx, app, nb.Pages[i].Submission.CanvasData)
		}
	}
}

func resolveNotebookFullSubmissionImages(ctx context.Context, app *appcontext.Context, items []domain.NotebookSubmissionFull) {
	for i := range items {
		items[i].CanvasData = datauri.Resolve(ctx, app, items[i].CanvasData)
	}
}
