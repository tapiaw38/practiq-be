package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

func schoolForLifecycle(ctx context.Context, app *appcontext.Context, id string) (*domain.School, apperrors.ApplicationError) {
	school, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	return school, nil
}
