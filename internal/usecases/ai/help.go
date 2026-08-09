package ai

import (
	"context"
	"log"
	"math/rand"
	"strings"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

var mockResponses = map[string][]string{
	"hint": {
		"Piensa en los pasos básicos: ¿cuál es el primer paso para resolver este tipo de problema?",
		"Recuerda la fórmula fundamental. ¿Qué datos tienes disponibles?",
		"Intenta dividir el problema en partes más pequeñas. ¿Qué es lo que ya sabes?",
		"Observa el patrón: ¿has visto algo similar antes?",
	},
	"explanation": {
		"Vamos paso a paso: primero identificamos los datos, luego aplicamos la operación correcta.",
		"El concepto clave aquí es entender la relación entre los elementos del problema.",
		"Este tipo de ejercicio requiere que apliques las reglas básicas que aprendiste. Piensa en el proceso.",
		"Analiza cada parte del problema por separado, luego une las piezas.",
	},
	"similar_example": {
		"Por ejemplo, si tienes 3/4 + 1/4, primero verificas que los denominadores sean iguales, luego sumas los numeradores.",
		"Imagina que divides una pizza en 8 partes. Si comes 3/8 y luego 2/8, ¿cuánto comiste en total? ¡Así se resuelve!",
		"Piensa en esto: si 2 × 3 = 6, entonces 2 × 30 = 60. ¿Ves el patrón?",
		"Como ejemplo: para resolver 15 ÷ 3, puedes pensar: ¿cuántas veces cabe el 3 en el 15? La respuesta es 5.",
	},
}

type (
	HelpUsecase interface {
		Execute(context.Context, HelpInput) (*HelpOutput, apperrors.ApplicationError)
	}

	helpUsecase struct {
		contextFactory appcontext.Factory
	}

	HelpInput struct {
		StudentID      string
		ExerciseID     string `json:"exercise_id"`
		Question       string `json:"question"`
		HelpType       string `json:"help_type"`
		ConversationID string `json:"conversation_id"`
	}

	HelpOutput struct {
		Data HelpData `json:"data"`
	}
)

func NewHelpUsecase(contextFactory appcontext.Factory) HelpUsecase {
	return &helpUsecase{contextFactory: contextFactory}
}

func (u *helpUsecase) Execute(ctx context.Context, input HelpInput) (*HelpOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	if input.ExerciseID != "" {
		exercise, err := app.Repositories.Exercise.Get(ctx, input.ExerciseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AIHelpError, err)
		}
		if exercise != nil && exercise.TopicID != "" {
			topic, _ := app.Repositories.Topic.Get(ctx, exercise.TopicID)
			if topic == nil {
				return nil, apperrors.NewForbiddenError()
			}
			hasAccess, err := studentHasCourseAccess(ctx, app, input.StudentID, topic.CourseID)
			if err != nil {
				return nil, apperrors.NewApplicationError(mappings.AIHelpError, err)
			}
			if !hasAccess {
				return nil, apperrors.NewForbiddenError()
			}
		}
	}

	helpType := input.HelpType
	if helpType == "" {
		helpType = "hint"
	}

	response := u.getAIResponse(ctx, app, input, helpType)

	id, err := app.Repositories.AIConversation.CreateHelpRequest(ctx, domain.AIHelpRequest{
		StudentID:  input.StudentID,
		ExerciseID: input.ExerciseID,
		Question:   input.Question,
		AIResponse: response,
		HelpType:   helpType,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.AIHelpError, err)
	}

	if input.ConversationID != "" {
		conversation, err := app.Repositories.AIConversation.Get(ctx, input.ConversationID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.AIConversationGetError, err)
		}
		if conversation == nil || conversation.StudentID != input.StudentID {
			return nil, apperrors.NewForbiddenError()
		}
		u.persistMessages(ctx, app, input.ConversationID, input.Question, helpType, response)
	}

	return &HelpOutput{Data: toHelpOutputData(id, response, helpType)}, nil
}

func (u *helpUsecase) getAIResponse(ctx context.Context, app *appcontext.Context, input HelpInput, helpType string) string {
	profile, err := app.Repositories.UserProfile.Get(ctx, input.StudentID)
	if err != nil {
		log.Printf("[ai_help] warning: failed to get user profile student_id=%s err=%v", input.StudentID, err)
		return getMockResponse(helpType)
	}
	if profile == nil || profile.AssistantBaseURL == "" {
		log.Printf("[ai_help] warning: assistant not configured for student_id=%s", input.StudentID)
		return getMockResponse(helpType)
	}

	var exercise *domain.Exercise
	gradeName := ""
	if input.ExerciseID != "" {
		exercise, err = app.Repositories.Exercise.Get(ctx, input.ExerciseID)
		if err != nil {
			log.Printf("[ai_help] warning: failed to get exercise exercise_id=%s err=%v", input.ExerciseID, err)
		}
		// Get grade from exercise → topic → course
		if exercise != nil && exercise.TopicID != "" {
			if topic, _ := app.Repositories.Topic.Get(ctx, exercise.TopicID); topic != nil {
				if course, _ := app.Repositories.Course.Get(ctx, topic.CourseID); course != nil {
					gradeName = course.GradeName
				}
			}
		}
	}

	history, historyErr := app.Repositories.AIConversation.ListRecentHelpRequests(ctx, input.StudentID, input.ExerciseID, 3)
	if historyErr != nil {
		log.Printf("[ai_help] warning: failed to load exercise memory student_id=%s exercise_id=%s err=%v", input.StudentID, input.ExerciseID, historyErr)
	}
	prompt := buildHelpPrompt(helpType, input.Question, exercise, gradeName, history)

	cfg := assistant.Config{
		BaseURL: profile.AssistantBaseURL,
		APIKey:  profile.AssistantAPIKey,
	}

	aiResponse, err := app.Integrations.AssistantGateway.AskHelp(ctx, cfg, prompt)
	if err != nil {
		log.Printf("[ai_help] warning: assistant call failed student_id=%s err=%v, falling back to mock", input.StudentID, err)
		return getMockResponse(helpType)
	}

	return strings.TrimSpace(aiResponse)
}

func buildHelpPrompt(helpType, studentQuestion string, exercise *domain.Exercise, gradeName string, history []domain.AIHelpRequest) string {
	var sb strings.Builder

	if gradeName != "" {
		sb.WriteString("Eres un tutor de matemáticas amigable para estudiantes de ")
		sb.WriteString(gradeName)
		sb.WriteString(". ")
		sb.WriteString("Considera los contenidos de los documentos de ")
		sb.WriteString(gradeName)
		sb.WriteString(" para responder. ")
	} else {
		sb.WriteString("Eres un tutor de matemáticas amigable para estudiantes de primaria y secundaria. ")
	}
	sb.WriteString("Responde siempre en español, de forma clara y pedagógica.\n\n")

	if exercise != nil {
		sb.WriteString("Contexto del ejercicio:\n")
		sb.WriteString("- Enunciado: ")
		sb.WriteString(exercise.Question)
		sb.WriteString("\n")
		if exercise.Difficulty > 0 {
			sb.WriteString("- Dificultad: ")
			switch exercise.Difficulty {
			case 1:
				sb.WriteString("básica")
			case 2:
				sb.WriteString("intermedia")
			case 3:
				sb.WriteString("avanzada")
			default:
				sb.WriteString("variable")
			}
			sb.WriteString("\n")
		}
	}

	switch helpType {
	case "hint":
		sb.WriteString("\nEl estudiante pide una PISTA. ")
		sb.WriteString("Dale una pista útil que lo guíe hacia la solución SIN revelar la respuesta directamente. ")
		sb.WriteString("Ayúdalo a pensar en el proceso.\n")
	case "explanation":
		sb.WriteString("\nEl estudiante pide una EXPLICACIÓN detallada. ")
		if exercise != nil && exercise.Explanation != "" {
			sb.WriteString("Referencia para tu explicación: ")
			sb.WriteString(exercise.Explanation)
			sb.WriteString("\n")
		}
		if exercise != nil && exercise.CorrectAnswer != "" {
			sb.WriteString("La respuesta correcta es: ")
			sb.WriteString(exercise.CorrectAnswer)
			sb.WriteString("\n")
		}
		sb.WriteString("Explica paso a paso cómo resolver este tipo de problema.\n")
	case "similar_example":
		sb.WriteString("\nEl estudiante pide un EJEMPLO SIMILAR. ")
		if exercise != nil && exercise.CorrectAnswer != "" {
			sb.WriteString("Basándote en el ejercicio original (respuesta: ")
			sb.WriteString(exercise.CorrectAnswer)
			sb.WriteString("), ")
		}
		sb.WriteString("crea un ejemplo similar pero con números diferentes y resuélvelo paso a paso.\n")
	default:
		sb.WriteString("\nAyuda al estudiante con su consulta de forma pedagógica.\n")
	}

	if studentQuestion != "" {
		sb.WriteString("\nPregunta del estudiante: ")
		sb.WriteString(studentQuestion)
		sb.WriteString("\n")
	}

	if len(history) > 0 {
		sb.WriteString("\nMemoria reciente de este ejercicio (no repitas ayuda ya dada):\n")
		for i := len(history) - 1; i >= 0; i-- {
			item := history[i]
			sb.WriteString("- Estudiante: ")
			sb.WriteString(item.Question)
			sb.WriteString("\n- Tutor: ")
			sb.WriteString(item.AIResponse)
			sb.WriteString("\n")
		}
	}

	sb.WriteString("\nResponde de forma concisa (máximo 3-4 oraciones para pistas, un poco más para explicaciones).")

	return sb.String()
}

func getMockResponse(helpType string) string {
	responses, ok := mockResponses[helpType]
	if !ok {
		responses = mockResponses["hint"]
	}
	return responses[rand.Intn(len(responses))]
}

func (u *helpUsecase) persistMessages(ctx context.Context, app *appcontext.Context, conversationID, question, helpType, aiResponse string) {
	studentContent := question
	if studentContent == "" {
		switch helpType {
		case "hint":
			studentContent = "Dame una pista"
		case "explanation":
			studentContent = "Explícame cómo resolver esto"
		case "similar_example":
			studentContent = "Muéstrame un ejemplo similar"
		default:
			studentContent = "Necesito ayuda"
		}
	}

	_, err := app.Repositories.AIConversation.AddMessage(ctx, domain.AIMessage{
		ConversationID: conversationID,
		Sender:         "student",
		MessageType:    "text",
		Content:        studentContent,
	})
	if err != nil {
		log.Printf("[ai_help] warning: failed to persist student message conversation_id=%s err=%v", conversationID, err)
	}

	_, err = app.Repositories.AIConversation.AddMessage(ctx, domain.AIMessage{
		ConversationID: conversationID,
		Sender:         "ai",
		MessageType:    "text",
		Content:        aiResponse,
	})
	if err != nil {
		log.Printf("[ai_help] warning: failed to persist ai message conversation_id=%s err=%v", conversationID, err)
	}
}

func studentHasCourseAccess(ctx context.Context, app *appcontext.Context, studentID, courseID string) (bool, error) {
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{StudentID: studentID})
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
