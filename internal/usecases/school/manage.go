package school

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
	"github.com/tapiaw38/practiq-be/internal/usecases/subscription"
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
		List(ctx context.Context, bearerToken string) (*SchoolsOutput, apperrors.ApplicationError)
		Create(ctx context.Context, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError)
		Update(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in SchoolInput) (*SchoolOutput, apperrors.ApplicationError)
		Close(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in CloseInput) (*SchoolOutput, apperrors.ApplicationError)
		Reopen(ctx context.Context, isSuperAdmin bool, id string) (*SchoolOutput, apperrors.ApplicationError)
		Archive(ctx context.Context, isSuperAdmin bool, id, bearerToken string) (*ArchiveOutput, apperrors.ApplicationError)
		// Mine is what the asking user belongs to, for the school selector.
		Mine(ctx context.Context, requesterID, bearerToken string) (*SchoolsOutput, apperrors.ApplicationError)
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

	CloseInput struct {
		ConfirmName string `json:"confirm_name"`
		Reason      string `json:"reason"`
	}

	SchoolData struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Kind    string `json:"kind"`
		Billing string `json:"billing"`
		Status  string `json:"status"`
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

	ArchiveCourseData struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		GradeName   string `json:"grade_name"`
		SubjectName string `json:"subject_name"`
	}
	ArchiveData struct {
		School  SchoolData          `json:"school"`
		Members []MemberData        `json:"members"`
		Courses []ArchiveCourseData `json:"courses"`
	}
	ArchiveOutput struct {
		Data ArchiveData `json:"data"`
	}
)

func NewManageUsecase(contextFactory appcontext.Factory) ManageUsecase {
	return &manageUsecase{contextFactory: contextFactory}
}

func (u *manageUsecase) List(ctx context.Context, bearerToken string) (*SchoolsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	schools, err := app.Repositories.School.List(ctx)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	schools, appErr := resolvePersonalSchoolNames(ctx, app, bearerToken, schools)
	if appErr != nil {
		return nil, appErr
	}

	data := make([]SchoolData, 0, len(schools))
	for _, s := range schools {
		data = append(data, toSchoolData(s, ""))
	}
	return &SchoolsOutput{Data: data}, nil
}

func (u *manageUsecase) Mine(ctx context.Context, requesterID, bearerToken string) (*SchoolsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	members, err := app.Repositories.School.ListForUser(ctx, requesterID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	schools := make([]domain.School, 0, len(members))
	roles := make(map[string]string, len(members))
	for _, member := range members {
		if !member.Active {
			continue
		}
		s, err := app.Repositories.School.Get(ctx, member.SchoolID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
		}
		if s != nil {
			schools = append(schools, *s)
			roles[s.ID] = member.Role
		}
	}
	schools, appErr := resolvePersonalSchoolNames(ctx, app, bearerToken, schools)
	if appErr != nil {
		return nil, appErr
	}

	data := make([]SchoolData, 0, len(schools))
	for _, school := range schools {
		data = append(data, toSchoolData(school, roles[school.ID]))
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

// Close preserves academic and billing history but removes every member from
// the active scope. Only a platform superadmin can close or reopen a school.
func (u *manageUsecase) Close(ctx context.Context, requesterID string, isSuperAdmin bool, id string, in CloseInput) (*SchoolOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	if school.Status == domain.SchoolStatusClosed {
		return u.read(ctx, app, id, "")
	}
	if strings.TrimSpace(in.ConfirmName) != school.Name {
		return nil, apperrors.NewBadRequestError("confirmation must match the school name")
	}
	if err := app.Repositories.School.Close(ctx, id, requesterID, strings.TrimSpace(in.Reason)); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	// Existing codes must become unusable with the school. RequireActive in
	// redeem is the second guard if a concurrent request races this update.
	if err := app.Repositories.StudentInvitation.RevokeForSchool(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.InvitationRevokeError, err)
	}
	return u.read(ctx, app, id, "")
}

func (u *manageUsecase) Reopen(ctx context.Context, isSuperAdmin bool, id string) (*SchoolOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	if err := app.Repositories.School.Reopen(ctx, id); err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	return u.read(ctx, app, id, "")
}

func (u *manageUsecase) Archive(ctx context.Context, isSuperAdmin bool, id, bearerToken string) (*ArchiveOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, err := app.Repositories.School.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	if school == nil {
		return nil, apperrors.NewNotFoundError("school not found")
	}
	members, appErr := u.ListMembers(ctx, "", true, id, bearerToken)
	if appErr != nil {
		return nil, appErr
	}
	courses, err := app.Repositories.Course.ListArchive(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseListError, err)
	}
	data := make([]ArchiveCourseData, 0, len(courses))
	for _, course := range courses {
		data = append(data, ArchiveCourseData{ID: course.ID, Title: course.Title, GradeName: course.GradeName, SubjectName: course.SubjectName})
	}
	return &ArchiveOutput{Data: ArchiveData{School: toSchoolData(*school, ""), Members: members.Data, Courses: data}}, nil
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

	// The same limit every other way in goes through. This one was open: the
	// "Usuarios" panel adds a student straight to the school, so without it a
	// teacher could pass their plan from the one screen built for it.
	if role == domain.SchoolRoleStudent {
		if appErr := subscription.EnsureCanAddStudent(ctx, app, schoolID, requesterID, in.UserID); appErr != nil {
			return appErr
		}
	}

	// A person joins a school by their Practiq profile, which exists once they
	// have signed in. Without this the insert failed on a foreign key and
	// surfaced as "failed to resolve the school" — a 500 about the wrong thing,
	// which sent an admin looking at the school instead of at the person they
	// were adding.
	profile, err := app.Repositories.UserProfile.Get(ctx, in.UserID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}
	if profile == nil {
		return apperrors.NewNotFoundError("that person has no Practiq profile yet — they have to sign in once before joining a school")
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
	// Removing the last admin would leave an active school nobody can
	// administer. This applies to a superadmin too: they can close/reopen, not
	// accidentally turn an operating school ownerless.
	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}
	for _, member := range members {
		if member.UserID == userID && member.Role == domain.SchoolRoleAdmin && member.Active {
			admins, err := app.Repositories.School.CountActiveAdmins(ctx, schoolID)
			if err != nil {
				return apperrors.NewApplicationError(mappings.SchoolLookupError, err)
			}
			if admins <= 1 {
				return apperrors.NewBadRequestError("an active school needs at least one admin")
			}
			break
		}
	}
	// Keep a normal school admin from removing themself even if another admin
	// exists; use a handover instead.
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

	if !isSuperAdmin {
		if appErr := EnsureAdministers(ctx, app, requesterID, false, schoolID); appErr != nil {
			return nil, appErr
		}
	}

	members, err := app.Repositories.School.ListMembers(ctx, schoolID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SchoolLookupError, err)
	}

	ids := make([]string, 0, len(members))
	for _, member := range members {
		ids = append(ids, member.UserID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if appErr != nil {
		return nil, appErr
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
	return SchoolData{ID: s.ID, Name: s.Name, Kind: s.Kind, Billing: s.Billing, Status: s.Status, Role: role}
}

// resolvePersonalSchoolNames fixes display names produced by the initial SQL
// migration. Auth owns legal names, so the database could only use created_by
// there. Keep an owner-renamed school intact; only replace known legacy names.
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
