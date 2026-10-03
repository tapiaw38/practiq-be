package course

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func mismatch(kind, id, rowSchool, selectedSchool string) apperrors.ApplicationError {
	return apperrors.NewApplicationError(
		mappings.ErrorDetails{
			InternalCode: "common:bad-request",
			StatusCode:   400,
			Message:      kind + " does not belong to active school",
		},
		fmt.Errorf("%s %s belongs to school %q, request selected %q", kind, id, rowSchool, selectedSchool),
	)
}

type (
	CreateUsecase interface {
		Execute(ctx context.Context, canCreate bool, teacherID, schoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		GradeID     string `json:"grade_id" binding:"required"`
		SubjectID   string `json:"subject_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description"`
		Level       string `json:"level"`
		Subject     string `json:"subject"`
	}

	CreateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, canCreate bool, teacherID, schoolID string, isSuperAdmin bool, in CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !canCreate {
		return nil, apperrors.NewForbiddenError()
	}
	if strings.TrimSpace(in.GradeID) == "" {
		return nil, apperrors.NewBadRequestError("grade_id is required")
	}
	if strings.TrimSpace(in.SubjectID) == "" {
		return nil, apperrors.NewBadRequestError("subject_id is required")
	}

	if strings.TrimSpace(schoolID) == "" {
		return nil, apperrors.NewBadRequestError("seleccioná una escuela antes de crear un curso")
	}
	grade, err := app.Repositories.Grade.Get(ctx, in.GradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewNotFoundError("grade not found")
	}
	if grade.SchoolID != schoolID {
		return nil, mismatch("grade", in.GradeID, grade.SchoolID, schoolID)
	}
	subjectRow, err := app.Repositories.Subject.Get(ctx, in.SubjectID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectGetError, err)
	}
	if subjectRow == nil {
		return nil, apperrors.NewNotFoundError("subject not found")
	}
	if subjectRow.SchoolID != schoolID {
		return nil, mismatch("subject", in.SubjectID, subjectRow.SchoolID, schoolID)
	}
	if appErr := school.EnsureAdministers(ctx, app, teacherID, isSuperAdmin, schoolID); appErr != nil {
		return nil, appErr
	}

	id, err := app.Repositories.Course.Create(ctx, domain.Course{
		TeacherID:   teacherID,
		GradeID:     in.GradeID,
		SubjectID:   in.SubjectID,
		Title:       in.Title,
		Description: in.Description,
		Level:       in.Level,
		Subject:     in.Subject,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
	}

	rollbackCourse := func() {
		if err := app.Repositories.Course.Delete(ctx, id); err != nil {
			log.Printf("[course_create] could not roll back course_id=%s err=%v", id, err)
		}
	}

	defaultStrategy, err := app.Repositories.LearningStrategy.GetByCode(ctx, "kumon")
	if err != nil {
		rollbackCourse()
		return nil, apperrors.NewApplicationError(mappings.LearningStrategyGetError, err)
	}
	if defaultStrategy != nil {
		if _, err := app.Repositories.LearningStrategy.AssignToCourse(ctx, domain.CourseLearningStrategy{
			CourseID:   id,
			StrategyID: defaultStrategy.ID,
			IsDefault:  true,
			Config:     "{}",
		}); err != nil {
			rollbackCourse()
			return nil, apperrors.NewApplicationError(mappings.LearningStrategyAssignError, err)
		}
	}

	c, err := app.Repositories.Course.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}

	return &CreateOutput{Data: toCourseData(*c)}, nil
}
