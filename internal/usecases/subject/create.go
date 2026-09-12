package subject

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
		Execute(context.Context, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string
		Description string
		CreatedBy   string
		SchoolID    string
	}

	CreateOutput struct {
		Data SubjectData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// Resolved here rather than taken from the request: a caller must not be
	// able to file a grade or subject under a school they do not belong to.
	schoolID, appErr := school.OwnedSchoolIDSelected(ctx, app, input.CreatedBy, input.SchoolID)
	if appErr != nil {
		return nil, appErr
	}
	id, err := app.Repositories.Subject.Create(ctx, domain.Subject{
		SchoolID:    schoolID,
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   input.CreatedBy,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectCreateError, err)
	}
	subject, err := app.Repositories.Subject.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectGetError, err)
	}
	if subject == nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectNotFoundError, nil)
	}
	return &CreateOutput{Data: toSubjectData(*subject)}, nil
}
