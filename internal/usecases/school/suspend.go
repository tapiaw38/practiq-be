package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	SuspendUsecase interface {
		Execute(ctx context.Context, isSuperAdmin bool, id string) (*SuspendOutput, apperrors.ApplicationError)
	}

	suspendUsecase struct {
		contextFactory appcontext.Factory
	}

	SuspendOutput struct {
		Data SchoolData `json:"data"`
	}
)

func NewSuspendUsecase(contextFactory appcontext.Factory) SuspendUsecase {
	return &suspendUsecase{contextFactory: contextFactory}
}

func (u *suspendUsecase) Execute(ctx context.Context, isSuperAdmin bool, id string) (*SuspendOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, appErr := schoolForLifecycle(ctx, app, id)
	if appErr != nil {
		return nil, appErr
	}
	if school.Status == domain.SchoolStatusClosed {
		return nil, apperrors.NewBadRequestError("a closed school must be reopened, not suspended")
	}
	if err := app.Repositories.School.Suspend(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	data, appErr := readSchool(ctx, app, id, "")
	if appErr != nil {
		return nil, appErr
	}
	return &SuspendOutput{Data: data}, nil
}
