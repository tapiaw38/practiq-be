package school

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	// ManageUsecase covers what an operator does with institutions and what an
	// institution's admin does with its people.
	//
	// Personal schools are not created here: they appear when a teacher signs
	// up. An institution exists because somebody agreed to invoice it, which is
	// a conversation, not a form a teacher fills in.
	ManageUsecase interface {
		// List and Create are the platform superadmin's.
		List(ctx context.Context) (*SchoolsOutput, apperrors.ApplicationError)
		Create(ctx context.Context, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError)
		Update(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError)
		// Mine is what the asking user belongs to, for the school selector.
		Mine(ctx context.Context, requesterID string) (*SchoolsOutput, apperrors.ApplicationError)
		// AddMember is the superadmin assigning an institution's admin, and
		// that admin adding teachers and students.
		AddMember(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string, in MemberInput) apperrors.ApplicationError
		RemoveMember(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, userID string) apperrors.ApplicationError
		ListMembers(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, bearerToken string) (*MembersOutput, apperrors.ApplicationError)
	}

	manageUsecase struct {
		contextFactory appcontext.Factory
	}

	SchoolInput struct {
		Name string `json:"name"`
		// Kind and Billing are only meaningful to a superadmin creating an
		// institution; a personal school is always subscription-billed.
		Kind    string `json:"kind"`
		Billing string `json:"billing"`
	}

	MemberInput struct {
		UserID string `json:"user_id"`
		Role   string `json:"role"`
	}

	SchoolData struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Billing string `json:"billing"`
		// Role is the asking user's role in it, empty when they are only
		// looking as a superadmin.
		Role string `json:"role,omitempty"`
	}

	MemberData struct {
		UserID string `json:"user_id"`
		Name   string `json:"name"`
		Email  string `json:"email"`
		Role   string `json:"role"`
		Active bool   `json:"active"`
	}

	SchoolsOutput struct {
		Data []SchoolData `json:"data"`
	}

	SchoolOutput struct {
		Data SchoolData `json:"data"`
	}

	MembersOutput struct {
		Data []MemberData `json:"data"`
	}
)

func NewManageUsecase(contextFactory appcontext.Factory) ManageUsecase {
	return &manageUsecase{contextFactory: contextFactory}
}

func (u *manageUsecase) List(ctx context.Context) (*SchoolsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	schools, err := app.Repositories.School.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	data := make([]SchoolData, 0, len(schools))
	for _, s := range schools {
		data = append(data, toSchoolData(s, ""))
	}
	return &SchoolsOutput{Data: data}, nil
}

func (u *manageUsecase) Mine(ctx context.Context, requesterID string) (*SchoolsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	members, err := app.Repositories.School.ListForUser(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	data := make([]SchoolData, 0, len(members))
	for _, member := range members {
		if !member.Active {
			continue
		}
		s, err := app.Repositories.School.Get(ctx, member.SchoolID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
		if s != nil {
			data = append(data, toSchoolData(*s, member.Role))
		}
	}
	return &SchoolsOutput{Data: data}, nil
}

func (u *manageUsecase) Create(ctx context.Context, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperrors.NewBadRequestError("a school needs a name")
	}

	kind := in.Kind
	if kind == "" {
		kind = domain.SchoolKindInstitution
	}
	billing := in.Billing
	if billing == "" {
		// An institution created by hand is invoiced by hand. Defaulting to
		// subscription would cap it at the free plan on its first student.
		billing = domain.SchoolBillingDirect
	}

	id, err := app.Repositories.School.Create(ctx, domain.School{
		Name: name, Kind: kind, Billing: billing,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return u.read(ctx, app, id, "")
}

func (u *manageUsecase) Update(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, id); appErr != nil {
		return nil, appErr
	}

	update := domain.School{Name: strings.TrimSpace(in.Name)}
	if isSuperAdmin {
		// What a school costs and what it allows are the operator's to change.
		// An admin renaming their own school must not be able to move it onto
		// direct billing and stop paying.
		update.Kind = in.Kind
		update.Billing = in.Billing
	}

	if err := app.Repositories.School.Update(ctx, id, update); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return u.read(ctx, app, id, "")
}

func (u *manageUsecase) AddMember(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID string, in MemberInput) apperrors.ApplicationError {
	app := u.contextFactory()

	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, schoolID); appErr != nil {
		return appErr
	}
	if strings.TrimSpace(in.UserID) == "" {
		return apperrors.NewBadRequestError("a member needs a user")
	}

	role := in.Role
	switch role {
	case domain.SchoolRoleAdmin, domain.SchoolRoleTeacher, domain.SchoolRoleStudent:
	default:
		return apperrors.NewBadRequestError("role must be admin, teacher or student")
	}

	school, err := app.Repositories.School.Get(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return apperrors.NewNotFoundError("school not found")
	}
	// A personal school is one teacher and their students. Letting its owner
	// add teachers would turn it into an institution without anyone agreeing to
	// invoice one.
	if school.Kind == domain.SchoolKindPersonal && role != domain.SchoolRoleStudent {
		return apperrors.NewForbiddenError()
	}

	if err := app.Repositories.School.AddMember(ctx, domain.SchoolMember{
		SchoolID: schoolID, UserID: in.UserID, Role: role,
	}); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}

func (u *manageUsecase) RemoveMember(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, userID string) apperrors.ApplicationError {
	app := u.contextFactory()

	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, schoolID); appErr != nil {
		return appErr
	}
	// Removing the last admin would leave a school nobody can administer, and
	// only a superadmin could put one back.
	if !isSuperAdmin && userID == requesterID {
		return apperrors.NewBadRequestError("you cannot remove yourself from a school you administer")
	}

	if err := app.Repositories.School.RemoveMember(ctx, schoolID, userID); err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return nil
}

func (u *manageUsecase) ListMembers(ctx context.Context, requesterID string, isSuperAdmin bool, schoolID, bearerToken string) (*MembersOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if appErr := EnsureAdministers(ctx, app, requesterID, isSuperAdmin, schoolID); appErr != nil {
		return nil, appErr
	}

	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.UserID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	data := make([]MemberData, 0, len(members))
	for _, m := range members {
		info := names[m.UserID]
		data = append(data, MemberData{
			UserID: m.UserID,
			Name:   identity.FullName(info, m.UserID),
			Email:  info.Email,
			Role:   m.Role,
			Active: m.Active,
		})
	}
	return &MembersOutput{Data: data}, nil
}

func (u *manageUsecase) read(ctx context.Context, app *appcontext.Context, id, role string) (*SchoolOutput, apperrors.ApplicationError) {
	s, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if s == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	return &SchoolOutput{Data: toSchoolData(*s, role)}, nil
}

func toSchoolData(s domain.School, role string) SchoolData {
	return SchoolData{ID: s.ID, Name: s.Name, Kind: s.Kind, Billing: s.Billing, Role: role}
}
