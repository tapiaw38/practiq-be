package notebook

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	UpdateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError)
	}

	UpdateInput struct {
		Title       string
		Description string
		TopicID     string
	}

	UpdateOutput struct {
		Data NotebookData `json:"data"`
	}

	updateUsecase struct{ contextFactory appcontext.Factory }
)

func NewUpdateUsecase(contextFactory appcontext.Factory) UpdateUsecase {
	return &updateUsecase{contextFactory: contextFactory}
}

func (u *updateUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string, input UpdateInput) (*UpdateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Verify notebook exists and check ownership
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	if nb == nil {
		return nil, apperrors.NewNotFoundError("notebook not found")
	}

	if _, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, nb.CourseID); appErr != nil {
		return nil, appErr
	}
	if appErr := ensureTopicBelongsToCourse(ctx, app, input.TopicID, nb.CourseID); appErr != nil {
		return nil, appErr
	}

	if err := app.Repositories.Notebook.Update(ctx, id, domain.Notebook{
		Title:       input.Title,
		Description: input.Description,
		TopicID:     input.TopicID,
	}); err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}

	nb, err = app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}

	return &UpdateOutput{Data: toNotebookData(nb)}, nil
}
