package notebook

import (
	"context"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	UpdatePageUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input UpdatePageInput) apperrors.ApplicationError
	}

	UpdatePageInput struct {
		PageID            string
		Title             string
		ContentType       string
		ContentData       string
		Instructions      string
		StatementText     *string
		StatementVerified *bool
	}

	updatePageUsecase struct{ contextFactory appcontext.Factory }
)

func NewUpdatePageUsecase(contextFactory appcontext.Factory) UpdatePageUsecase {
	return &updatePageUsecase{contextFactory: contextFactory}
}

func (u *updatePageUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, input UpdatePageInput) apperrors.ApplicationError {
	app := u.contextFactory()
	page, err := app.Repositories.Notebook.GetPage(ctx, input.PageID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if page == nil {
		return apperrors.NewNotFoundError("notebook page not found")
	}
	notebook, err := app.Repositories.Notebook.Get(ctx, page.NotebookID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if notebook == nil {
		return apperrors.NewNotFoundError("notebook not found")
	}
	if !isSuperAdmin && notebook.TeacherID != requesterID {
		return apperrors.NewForbiddenError()
	}

	contentData := input.ContentData
	if isLikelyImageData(contentData) && app.ImageStorage != nil {
		if uploaded, err := app.ImageStorage.UploadDataURI(ctx, "notebook", notebook.TeacherID, contentData); err == nil {
			contentData = uploaded
		} else {
			log.Printf("[image_storage] notebook page update upload failed page_id=%s err=%v", input.PageID, err)
		}
	}
	statementText := page.StatementText
	statementVerified := page.StatementVerified

	if input.StatementText != nil {
		statementText = strings.TrimSpace(*input.StatementText)
		statementVerified = true
	}
	if input.StatementVerified != nil {
		statementVerified = *input.StatementVerified
	}

	if contentData != page.ContentData {
		statementVerified = false
		statementText = transcribePageStatement(ctx, app, notebook.TeacherID, contentData, domain.NotebookPage{
			PageNumber:   page.PageNumber,
			Title:        input.Title,
			Instructions: input.Instructions,
		})
	}

	if err := app.Repositories.Notebook.UpdatePage(ctx, domain.NotebookPage{
		ID:                input.PageID,
		Title:             input.Title,
		ContentType:       input.ContentType,
		ContentData:       contentData,
		StatementText:     statementText,
		StatementVerified: statementVerified,
		Instructions:      input.Instructions,
	}); err != nil {
		return apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}
	return nil
}
