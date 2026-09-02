package enrollment

import (
	"context"

	enrollmentRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/enrollment"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

type (
	ListStudentsUsecase interface {
		Execute(context.Context, ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError)
	}

	listStudentsUsecase struct {
		contextFactory appcontext.Factory
	}

	ListStudentsInput struct {
		CourseID     string
		RequesterID  string
		IsSuperAdmin bool
		Limit        int
		Offset       int
		BearerToken  string
	}

	ListStudentsOutput struct {
		Data []StudentData `json:"data"`
	}
)

func NewListStudentsUsecase(contextFactory appcontext.Factory) ListStudentsUsecase {
	return &listStudentsUsecase{contextFactory: contextFactory}
}

func (u *listStudentsUsecase) Execute(ctx context.Context, input ListStudentsInput) (*ListStudentsOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// El padrón de un curso es dato de terceros: lo ve quien dicta ese curso.
	if !input.IsSuperAdmin {
		course, err := app.Repositories.Course.Get(ctx, input.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if course == nil {
			return nil, apperrors.NewNotFoundError("course not found")
		}
		if course.TeacherID != input.RequesterID {
			return nil, apperrors.NewForbiddenError()
		}
	}

	filter := enrollmentRepo.ListFilter{
		CourseID: input.CourseID,
		Limit:    input.Limit,
		Offset:   input.Offset,
	}

	students, err := app.Repositories.Enrollment.ListStudents(ctx, filter)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.EnrollmentListError, err)
	}

	ids := make([]string, 0, len(students))
	for _, s := range students {
		ids = append(ids, s.ID)
	}
	names, err := identity.Names(ctx, app.Integrations.AuthAPI, input.BearerToken, ids)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProfileGetError, err)
	}

	var data []StudentData
	for _, s := range students {
		info := names[s.ID]
		data = append(data, toStudentData(s, identity.FullName(info, s.ID), info.Email))
	}
	if data == nil {
		data = []StudentData{}
	}

	return &ListStudentsOutput{Data: data}, nil
}
