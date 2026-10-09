package ai

import (
	"context"
	"log"
	"strings"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

var defaultCuriosities = []string{
	"¿Sabías que las matemáticas se usan en los videojuegos?",
	"Los números nos ayudan a entender el mundo",
	"Aprender es como entrenar un superpoder",
	"Tu cerebro puede hacer cosas increíbles",
	"Cada problema resuelto te hace más fuerte",
}

type (
	GenerateCuriositiesUsecase interface {
		Execute(ctx context.Context, userID string, isSuperAdmin bool, in GenerateCuriositiesInput) (*GenerateCuriositiesOutput, apperrors.ApplicationError)
	}

	generateCuriositiesUsecase struct {
		contextFactory appcontext.Factory
	}

	GenerateCuriositiesInput struct {
		CourseID string `json:"course_id" binding:"required"`
	}

	GenerateCuriositiesOutput struct {
		Data CuriositiesData `json:"data"`
	}

	CuriositiesData struct {
		CourseID    string   `json:"course_id"`
		Curiosities []string `json:"curiosities"`
	}
)

func NewGenerateCuriositiesUsecase(contextFactory appcontext.Factory) GenerateCuriositiesUsecase {
	return &generateCuriositiesUsecase{contextFactory: contextFactory}
}

func (u *generateCuriositiesUsecase) Execute(ctx context.Context, userID string, isSuperAdmin bool, in GenerateCuriositiesInput) (*GenerateCuriositiesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	course, err := app.Repositories.Course.Get(ctx, in.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}

	_, manageErr := school.EnsureCanManageCourse(ctx, app, userID, isSuperAdmin, in.CourseID)
	manages := manageErr == nil
	if !manages {
		hasAccess, err := userHasCourseAccess(ctx, app, userID, in.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

	cached, err := app.Repositories.CourseCuriosities.Get(ctx, in.CourseID)
	if err != nil {
		log.Printf("[ai_curiosities] warning: failed to get cached curiosities course_id=%s err=%v", in.CourseID, err)
	}
	if cached != nil && len(cached.Curiosities) > 0 && !containsTechnicalFallback(cached.Curiosities) {
		return &GenerateCuriositiesOutput{
			Data: CuriositiesData{
				CourseID:    in.CourseID,
				Curiosities: cached.Curiosities,
			},
		}, nil
	}
	if cached != nil && len(cached.Curiosities) > 0 {
		log.Printf("[ai_curiosities] warning: discarding technical fallback from cache course_id=%s", in.CourseID)
	}

	cfg := assistantcfg.Resolve(ctx, app)
	if !app.Integrations.AssistantGateway.IsConfigured(cfg) {
		log.Print("[ai_curiosities] warning: the assistant is not configured for this platform")
		return u.fallbackResponse(in.CourseID), nil
	}

	subject := course.SubjectName
	if subject == "" {
		subject = course.Subject
	}
	topic := course.Title
	if course.Description != "" {
		topic = topic + " - " + course.Description
	}

	curiosities, err := app.Integrations.AssistantGateway.GenerateCourseCuriosities(ctx, cfg, subject, topic, course.GradeName, 8)
	if err != nil {
		log.Printf("[ai_curiosities] warning: AI generation failed course_id=%s err=%v, falling back to defaults", in.CourseID, err)
		return u.fallbackResponse(in.CourseID), nil
	}

	if len(curiosities) == 0 {
		return u.fallbackResponse(in.CourseID), nil
	}

	if manages {
		if err := app.Repositories.CourseCuriosities.Upsert(ctx, domain.CourseCuriosities{
			CourseID:    in.CourseID,
			Curiosities: curiosities,
		}); err != nil {
			log.Printf("[ai_curiosities] warning: failed to cache curiosities course_id=%s err=%v", in.CourseID, err)
		}
	}

	return &GenerateCuriositiesOutput{
		Data: CuriositiesData{
			CourseID:    in.CourseID,
			Curiosities: curiosities,
		},
	}, nil
}

func (u *generateCuriositiesUsecase) fallbackResponse(courseID string) *GenerateCuriositiesOutput {
	return &GenerateCuriositiesOutput{
		Data: CuriositiesData{
			CourseID:    courseID,
			Curiosities: defaultCuriosities,
		},
	}
}

func containsTechnicalFallback(curiosities []string) bool {
	for _, curiosity := range curiosities {
		text := strings.ToLower(curiosity)
		if strings.Contains(text, "problemas técnicos") || strings.Contains(text, "intentarlo de nuevo en unos minutos") {
			return true
		}
	}
	return false
}

func userHasCourseAccess(ctx context.Context, app *appcontext.Context, userID, courseID string) (bool, error) {
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{StudentID: userID})
	if err != nil {
		return false, err
	}
	for _, course := range courses {
		if course.ID == courseID {
			return true, nil
		}
	}
	return false, nil
}
