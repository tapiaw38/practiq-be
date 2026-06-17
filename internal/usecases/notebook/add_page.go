package notebook

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	AddPageUsecase interface {
		Execute(ctx context.Context, input AddPageInput) (*AddPageOutput, error)
	}

	AddPageInput struct {
		NotebookID   string
		PageNumber   int
		Title        string
		ContentType  string
		ContentData  string
		Instructions string
	}

	AddPageOutput struct {
		Data PageData `json:"data"`
	}

	addPageUsecase struct{ contextFactory appcontext.Factory }
)

func NewAddPageUsecase(contextFactory appcontext.Factory) AddPageUsecase {
	return &addPageUsecase{contextFactory: contextFactory}
}

func (u *addPageUsecase) Execute(ctx context.Context, input AddPageInput) (*AddPageOutput, error) {
	app := u.contextFactory()
	contentData := input.ContentData
	if isLikelyImageData(contentData) && app.ImageStorage != nil {
		notebook, err := app.Repositories.Notebook.Get(ctx, input.NotebookID)
		if err != nil {
			return nil, err
		}
		userID := "unknown"
		if notebook != nil {
			userID = notebook.TeacherID
		}
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", userID, contentData); err == nil {
			contentData = uploaded
		} else {
			log.Printf("[image_storage] notebook page upload failed notebook_id=%s err=%v", input.NotebookID, err)
		}
	}
	id, err := app.Repositories.Notebook.CreatePage(ctx, domain.NotebookPage{
		NotebookID:   input.NotebookID,
		PageNumber:   input.PageNumber,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  contentData,
		Instructions: input.Instructions,
	})
	if err != nil {
		return nil, err
	}
	return &AddPageOutput{Data: toPageData(id, input, contentData)}, nil
}
