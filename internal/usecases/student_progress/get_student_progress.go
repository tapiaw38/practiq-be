package studentprogress

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetStudentProgressUsecase interface {
	Execute(ctx context.Context, studentID string) (*ProgressListOutput, apperrors.ApplicationError)
}

type getStudentProgressUsecase struct {
	factory appcontext.Factory
}

func NewGetStudentProgressUsecase(factory appcontext.Factory) GetStudentProgressUsecase {
	return &getStudentProgressUsecase{factory: factory}
}

func (u *getStudentProgressUsecase) Execute(ctx context.Context, studentID string) (*ProgressListOutput, apperrors.ApplicationError) {
	app := u.factory()

	list, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p))
	}
	if data == nil {
		data = []ProgressData{}
	}

	return &ProgressListOutput{Data: data}, nil
}

type GetStudentCourseProgressUsecase interface {
	Execute(ctx context.Context, studentID, courseID string) (*ProgressListOutput, apperrors.ApplicationError)
}

type getStudentCourseProgressUsecase struct {
	factory appcontext.Factory
}

func NewGetStudentCourseProgressUsecase(factory appcontext.Factory) GetStudentCourseProgressUsecase {
	return &getStudentCourseProgressUsecase{factory: factory}
}

func (u *getStudentCourseProgressUsecase) Execute(ctx context.Context, studentID, courseID string) (*ProgressListOutput, apperrors.ApplicationError) {
	app := u.factory()

	list, err := app.Repositories.StudentProgress.ListByStudentAndCourse(ctx, studentID, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	var data []ProgressData
	for _, p := range list {
		data = append(data, toProgressData(p))
	}
	if data == nil {
		data = []ProgressData{}
	}

	return &ProgressListOutput{Data: data}, nil
}
