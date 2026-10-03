package subject

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/pgerr"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

const subjectNameConstraint = "idx_subjects_school_name"

type (
	CreateUsecase interface {
		Execute(ctx context.Context, createdBy, schoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	CreateOutput struct {
		Data SubjectData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, createdBy, requestedSchoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	schoolID, appErr := school.OwnedSchoolIDSelected(ctx, app, createdBy, isSuperAdmin, requestedSchoolID)
	if appErr != nil {
		return nil, appErr
	}
	id, err := app.Repositories.Subject.Create(ctx, domain.Subject{
		SchoolID:    schoolID,
		Name:        in.Name,
		Description: in.Description,
		CreatedBy:   createdBy,
	})
	if err != nil {
		if pgerr.IsUniqueViolation(err, subjectNameConstraint) {
			return nil, apperrors.NewConflictError("ya existe una materia con ese nombre en esta escuela")
		}
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
