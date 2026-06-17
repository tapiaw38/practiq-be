package notebook

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
)

type (
	UpdatePageUsecase interface {
		Execute(ctx context.Context, input UpdatePageInput) error
	}

	UpdatePageInput struct {
		PageID       string
		Title        string
		ContentType  string
		ContentData  string
		Instructions string
	}

	updatePageUsecase struct{ contextFactory appcontext.Factory }
)

func NewUpdatePageUsecase(contextFactory appcontext.Factory) UpdatePageUsecase {
	return &updatePageUsecase{contextFactory: contextFactory}
}

func (u *updatePageUsecase) Execute(ctx context.Context, input UpdatePageInput) error {
	app := u.contextFactory()
	contentData := input.ContentData
	if isLikelyImageData(contentData) && app.ImageStorage != nil {
		userID := "unknown"
		page, err := app.Repositories.Notebook.GetPage(ctx, input.PageID)
		if err != nil {
			return err
		}
		if page != nil {
			notebook, err := app.Repositories.Notebook.Get(ctx, page.NotebookID)
			if err != nil {
				return err
			}
			if notebook != nil {
				userID = notebook.TeacherID
			}
		}
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", userID, contentData); err == nil {
			contentData = uploaded
		} else {
			log.Printf("[image_storage] notebook page update upload failed page_id=%s err=%v", input.PageID, err)
		}
	}
	return app.Repositories.Notebook.UpdatePage(ctx, domain.NotebookPage{
		ID:           input.PageID,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  contentData,
		Instructions: input.Instructions,
	})
}
