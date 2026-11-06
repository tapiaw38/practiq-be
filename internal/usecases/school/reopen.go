package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ReopenUsecase interface {
		Execute(ctx context.Context, isSuperAdmin bool, id string) (*ReopenOutput, apperrors.ApplicationError)
	}

	reopenUsecase struct {
		contextFactory appcontext.Factory
	}

	ReopenOutput struct {
		Data SchoolData `json:"data"`
	}
)

func NewReopenUsecase(contextFactory appcontext.Factory) ReopenUsecase {
	return &reopenUsecase{contextFactory: contextFactory}
}

func (u *reopenUsecase) Execute(ctx context.Context, isSuperAdmin bool, id string) (*ReopenOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	if _, appErr := schoolForLifecycle(ctx, app, id); appErr != nil {
		return nil, appErr
	}
	if err := app.Repositories.School.Reopen(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	data, appErr := readSchool(ctx, app, id, "")
	if appErr != nil {
		return nil, appErr
	}
	return &ReopenOutput{Data: data}, nil
}
