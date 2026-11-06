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
	CreateUsecase interface {
		Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	CreateInput struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Level       int    `json:"level"`
		TopicID     string `json:"topic_id" binding:"required"`
	}

	CreateOutput struct {
		Data NotebookData `json:"data"`
	}

	createUsecase struct{ contextFactory appcontext.Factory }
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, courseID string, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	course, appErr := school.EnsureCanManageCourse(ctx, app, requesterID, isSuperAdmin, courseID)
	if appErr != nil {
		return nil, appErr
	}
	if appErr := ensureTopicBelongsToCourse(ctx, app, input.TopicID, course.ID); appErr != nil {
		return nil, appErr
	}

	owner := course.TeacherID
	if owner == "" {
		owner = requesterID
	}

	id, err := app.Repositories.Notebook.Create(ctx, domain.Notebook{
		CourseID:    courseID,
		TopicID:     input.TopicID,
		TeacherID:   owner,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookUpdateError, err)
	}
	nb, err := app.Repositories.Notebook.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.NotebookGetError, err)
	}
	resolveNotebookImages(ctx, app, nb)
	return &CreateOutput{Data: toNotebookData(nb)}, nil
}
