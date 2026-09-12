package notebook

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/datauri"
)

type (
	AddPageUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input AddPageInput) (*AddPageOutput, apperrors.ApplicationError)
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

func (u *addPageUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input AddPageInput) (*AddPageOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	notebook, err := app.Repositories.Notebook.Get(ctx, input.NotebookID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if notebook == nil {
		return nil, apperrors.NewNotFoundError("notebook not found")
	}
	if !isSuperAdmin && notebook.TeacherID != requesterID {
		return nil, apperrors.NewForbiddenError()
	}

	contentData := input.ContentData
	if isLikelyImageData(contentData) && app.ImageStorage != nil {
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", notebook.TeacherID, contentData); err == nil {
			contentData = uploaded
		} else {
			log.Printf("[image_storage] notebook page upload failed notebook_id=%s err=%v", input.NotebookID, err)
		}
	}
	page := domain.NotebookPage{
		NotebookID:   input.NotebookID,
		PageNumber:   input.PageNumber,
		Title:        input.Title,
		ContentType:  input.ContentType,
		ContentData:  contentData,
		Instructions: input.Instructions,
	}
	page.StatementText = transcribePageStatement(ctx, app, notebook.TeacherID, contentData, page)

	id, err := app.Repositories.Notebook.CreatePage(ctx, page)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}
	data := toPageData(id, input, contentData)
	data.ContentData = datauri.Resolve(ctx, app, data.ContentData)
	data.StatementText = page.StatementText
	return &AddPageOutput{Data: data}, nil
}
