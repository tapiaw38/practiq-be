package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

func readSchool(ctx context.Context, app *appcontext.Context, id, role string) (SchoolData, apperrors.ApplicationError) {
	school, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return SchoolData{}, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return SchoolData{}, apperrors.NewNotFoundError("school not found")
	}
	return toSchoolData(*school, role), nil
}

func toSchoolData(s domain.School, role string) SchoolData {
	return SchoolData{ID: s.ID, Name: s.Name, Kind: s.Kind, Billing: s.Billing, Status: s.Status, Role: role}
}

func resolvePersonalSchoolNames(ctx context.Context, app *appcontext.Context, bearerToken string, schools []domain.School) ([]domain.School, apperrors.ApplicationError) {
	ownerIDs := make([]string, 0, len(schools))
	for _, school := range schools {
		if school.Kind == domain.SchoolKindPersonal && legacyPersonalSchoolName(school) {
			ownerIDs = append(ownerIDs, school.CreatedBy)
		}
	}

	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ownerIDs)
	if appErr != nil {
		return nil, appErr
	}
	for i := range schools {
		if legacyPersonalSchoolName(schools[i]) {
			schools[i].Name = domain.PersonalSchoolName(identity.FullName(names[schools[i].CreatedBy], schools[i].CreatedBy))
		}
	}
	return schools, nil
}

func legacyPersonalSchoolName(school domain.School) bool {
	return school.Kind == domain.SchoolKindPersonal &&
		(school.Name == domain.PlaceholderSchoolName || school.Name == domain.PersonalSchoolName(school.CreatedBy))
}
