package school

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	ListUsecase interface {
		Execute(ctx context.Context, bearerToken string) (*ListOutput, apperrors.ApplicationError)
	}

	listUsecase struct {
		contextFactory appcontext.Factory
	}

	ListOutput struct {
		Data []SchoolData `json:"data"`
	}
)

func NewListUsecase(contextFactory appcontext.Factory) ListUsecase {
	return &listUsecase{contextFactory: contextFactory}
}

func (u *listUsecase) Execute(ctx context.Context, bearerToken string) (*ListOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	schools, err := app.Repositories.School.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	schools, appErr := resolvePersonalSchoolNames(ctx, app, bearerToken, schools)
	if appErr != nil {
		return nil, appErr
	}
	ownerIDs := make([]string, 0, len(schools))
	for _, school := range schools {
		if school.Kind == domain.SchoolKindPersonal && school.CreatedBy != "" {
			ownerIDs = append(ownerIDs, school.CreatedBy)
		}
	}
	owners, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ownerIDs)
	if appErr != nil {
		log.Printf("[school] owner lookup failed err=%v", appErr)
		owners = nil
	}
	data := make([]SchoolData, 0, len(schools))
	for _, school := range schools {
		item := toSchoolData(school, "")
		if school.Kind == domain.SchoolKindPersonal && school.CreatedBy != "" {
			info := owners[school.CreatedBy]
			item.Owner = &SchoolOwner{ID: school.CreatedBy, Name: identity.FullName(info, school.CreatedBy), Email: info.Email}
			item.Plan = ownerPlan(ctx, app, school.CreatedBy)
		}
		data = append(data, item)
	}
	return &ListOutput{Data: data}, nil
}

func ownerPlan(ctx context.Context, app *appcontext.Context, ownerID string) *SchoolPlan {
	entitlement, err := app.Integrations.Payments.GetEntitlement(ctx, ownerID)
	if err != nil {
		log.Printf("[payments] plan lookup failed owner_id=%s err=%v", ownerID, err)
		return nil
	}
	if entitlement == nil || !entitlement.Active {
		return &SchoolPlan{Name: domain.FreePlan.Name, MaxStudents: domain.FreePlan.MaxStudents}
	}
	planID := 0
	if entitlement.PlanID != nil {
		planID = *entitlement.PlanID
	}
	name, _ := entitlement.Metadata["name"].(string)
	plan := domain.PlanFromMetadata(planID, name, entitlement.Metadata)
	return &SchoolPlan{Name: plan.Name, MaxStudents: plan.MaxStudents, Active: true}
}
