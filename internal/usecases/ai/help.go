package ai

import (
	"context"
	"math/rand"

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

type HelpUsecase interface {
	Execute(context.Context, HelpInput) (*HelpOutput, apperrors.ApplicationError)
}

type helpUsecase struct {
	factory appcontext.Factory
}

type HelpInput struct {
	StudentID  string
	ExerciseID string `json:"exercise_id"`
	Question   string `json:"question"`
	HelpType   string `json:"help_type"`
}

func NewHelpUsecase(factory appcontext.Factory) HelpUsecase {
	return &helpUsecase{factory: factory}
}

func (u *helpUsecase) Execute(ctx context.Context, input HelpInput) (*HelpOutput, apperrors.ApplicationError) {
	app := u.factory()

	helpType := input.HelpType
	if helpType == "" {
		helpType = "hint"
	}

	responses, ok := mockResponses[helpType]
	if !ok {
		responses = mockResponses["hint"]
	}

	response := responses[rand.Intn(len(responses))]

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

	return &HelpOutput{Data: HelpData{
		ID:       id,
		Response: response,
		HelpType: helpType,
	}}, nil
}
