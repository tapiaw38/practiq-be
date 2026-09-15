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

// mismatch reports a grade or subject that belongs to another school.
//
// The ids travel as the original message so the log says which row and which
// two schools disagreed. Without them this refusal was indistinguishable in
// the log from the missing-header case above, which is what made a real
// 400 in production impossible to tell apart from a client that simply had
// no school selected.
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
		Execute(context.Context, bool, CreateInput) (*CreateOutput, apperrors.ApplicationError)
	}

	createUsecase struct {
		contextFactory appcontext.Factory
	}

	CreateInput struct {
		TeacherID    string
		SchoolID     string
		IsSuperAdmin bool
		GradeID      string `json:"grade_id"`
		SubjectID    string `json:"subject_id"`
		Title        string `json:"title"`
		Description  string `json:"description"`
		Level        string `json:"level"`
		Subject      string `json:"subject"`
	}

	CreateOutput struct {
		Data CourseData `json:"data"`
	}
)

func NewCreateUsecase(contextFactory appcontext.Factory) CreateUsecase {
	return &createUsecase{contextFactory: contextFactory}
}

func (u *createUsecase) Execute(ctx context.Context, canCreate bool, input CreateInput) (*CreateOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if !canCreate {
		return nil, apperrors.NewForbiddenError()
	}
	if strings.TrimSpace(input.GradeID) == "" {
		return nil, apperrors.NewBadRequestError("grade_id is required")
	}
	if strings.TrimSpace(input.SubjectID) == "" {
		return nil, apperrors.NewBadRequestError("subject_id is required")
	}
	// Three different failures used to answer with one message here, and the
	// one that actually happens most reads as the most misleading: with no
	// school selected, SchoolID is empty, every grade compares unequal to it,
	// and the caller is told their grade belongs elsewhere when the request
	// simply never said where "here" is.
	if strings.TrimSpace(input.SchoolID) == "" {
		return nil, apperrors.NewBadRequestError("seleccioná una escuela antes de crear un curso")
	}
	grade, err := app.Repositories.Grade.Get(ctx, input.GradeID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.GradeGetError, err)
	}
	if grade == nil {
		return nil, apperrors.NewNotFoundError("grade not found")
	}
	if grade.SchoolID != input.SchoolID {
		return nil, mismatch("grade", input.GradeID, grade.SchoolID, input.SchoolID)
	}
	subject, err := app.Repositories.Subject.Get(ctx, input.SubjectID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.SubjectGetError, err)
	}
	if subject == nil {
		return nil, apperrors.NewNotFoundError("subject not found")
	}
	if subject.SchoolID != input.SchoolID {
		return nil, mismatch("subject", input.SubjectID, subject.SchoolID, input.SchoolID)
	}
	if appErr := school.EnsureAdministers(ctx, app, input.TeacherID, input.IsSuperAdmin, input.SchoolID); appErr != nil {
		return nil, appErr
	}

	id, err := app.Repositories.Course.Create(ctx, domain.Course{
		TeacherID:   input.TeacherID,
		GradeID:     input.GradeID,
		SubjectID:   input.SubjectID,
		Title:       input.Title,
		Description: input.Description,
		Level:       input.Level,
		Subject:     input.Subject,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseCreateError, err)
	}

	// The repositories share a *sql.DB with no transaction plumbing, so failing
	// after the insert used to leave the course committed: the client retried
	// the "failed" creation and ended up with duplicates, while the first one
	// sat there without its default strategy. Undo the insert instead.
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
