package school

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	ArchiveUsecase interface {
		Execute(ctx context.Context, isSuperAdmin bool, id, bearerToken string) (*ArchiveOutput, apperrors.ApplicationError)
	}

	archiveUsecase struct {
		contextFactory appcontext.Factory
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

func NewArchiveUsecase(contextFactory appcontext.Factory) ArchiveUsecase {
	return &archiveUsecase{contextFactory: contextFactory}
}

func (u *archiveUsecase) Execute(ctx context.Context, isSuperAdmin bool, id, bearerToken string) (*ArchiveOutput, apperrors.ApplicationError) {
	if !isSuperAdmin {
		return nil, apperrors.NewForbiddenError()
	}
	app := u.contextFactory()
	school, appErr := schoolForLifecycle(ctx, app, id)
	if appErr != nil {
		return nil, appErr
	}
	members, appErr := listMembers(ctx, app, "", true, id, bearerToken)
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
	return &ArchiveOutput{Data: ArchiveData{School: toSchoolData(*school, ""), Members: members, Courses: data}}, nil
}
