package ai

import (
	"context"
	"log"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
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
		Execute(context.Context, GenerateCuriositiesInput) (*GenerateCuriositiesOutput, apperrors.ApplicationError)
	}

	generateCuriositiesUsecase struct {
		contextFactory appcontext.Factory
	}

	GenerateCuriositiesInput struct {
		UserID   string
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

func (u *generateCuriositiesUsecase) Execute(ctx context.Context, input GenerateCuriositiesInput) (*GenerateCuriositiesOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	course, err := app.Repositories.Course.Get(ctx, input.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
	}
	if course == nil {
		return nil, apperrors.NewNotFoundError("course not found")
	}

	if course.TeacherID != input.UserID {
		hasAccess, err := userHasCourseAccess(ctx, app, input.UserID, input.CourseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.CourseGetError, err)
		}
		if !hasAccess {
			return nil, apperrors.NewForbiddenError()
		}
	}

	// Check if course already has cached curiosities
	cached, err := app.Repositories.CourseCuriosities.Get(ctx, input.CourseID)
	if err != nil {
		log.Printf("[ai_curiosities] warning: failed to get cached curiosities course_id=%s err=%v", input.CourseID, err)
	}
	if cached != nil && len(cached.Curiosities) > 0 {
		return &GenerateCuriositiesOutput{
			Data: CuriositiesData{
				CourseID:    input.CourseID,
				Curiosities: cached.Curiosities,
			},
		}, nil
	}

	// Get user profile for AI config
	profile, err := app.Repositories.UserProfile.Get(ctx, input.UserID)
	if err != nil {
		log.Printf("[ai_curiosities] warning: failed to get user profile user_id=%s err=%v", input.UserID, err)
		return u.fallbackResponse(input.CourseID), nil
	}
	if profile == nil || profile.AssistantBaseURL == "" {
		log.Printf("[ai_curiosities] warning: assistant not configured for user_id=%s", input.UserID)
		return u.fallbackResponse(input.CourseID), nil
	}

	cfg := assistant.Config{
		BaseURL: profile.AssistantBaseURL,
		APIKey:  profile.AssistantAPIKey,
	}

	// Build topic string from course info
	subject := course.SubjectName
	if subject == "" {
		subject = course.Subject
	}
	topic := course.Title
	if course.Description != "" {
		topic = topic + " - " + course.Description
	}

	// Generate curiosities via AI
	curiosities, err := app.Integrations.AssistantGateway.GenerateCourseCuriosities(ctx, cfg, subject, topic, course.GradeName, 8)
	if err != nil {
		log.Printf("[ai_curiosities] warning: AI generation failed course_id=%s err=%v, falling back to defaults", input.CourseID, err)
		return u.fallbackResponse(input.CourseID), nil
	}

	if len(curiosities) == 0 {
		return u.fallbackResponse(input.CourseID), nil
	}

	if course.TeacherID == input.UserID {
		if err := app.Repositories.CourseCuriosities.Upsert(ctx, domain.CourseCuriosities{
			CourseID:    input.CourseID,
			Curiosities: curiosities,
		}); err != nil {
			log.Printf("[ai_curiosities] warning: failed to cache curiosities course_id=%s err=%v", input.CourseID, err)
		}
	}

	return &GenerateCuriositiesOutput{
		Data: CuriositiesData{
			CourseID:    input.CourseID,
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
