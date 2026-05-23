package practicesheet

import (
	"context"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	"github.com/tapiaw38/practiq-be/internal/platform/assistant"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/strategy"
)

type SubmitUsecase interface {
	Execute(context.Context, string, string, SubmitInput) (*SubmitOutput, apperrors.ApplicationError)
}

type submitUsecase struct {
	factory appcontext.Factory
}

type AttemptInput struct {
	ExerciseID       string `json:"exercise_id"`
	AnswerText       string `json:"answer_text"`
	CanvasData       string `json:"canvas_data"` // base64 PNG for canvas/handwritten exercises
	TimeSpentSeconds int    `json:"time_spent_seconds"`
	HintsUsed        int    `json:"hints_used"`
}

type SubmitInput struct {
	Attempts []AttemptInput `json:"attempts"`
}

func NewSubmitUsecase(factory appcontext.Factory) SubmitUsecase {
	return &submitUsecase{factory: factory}
}

func (u *submitUsecase) Execute(ctx context.Context, sheetID, studentID string, input SubmitInput) (*SubmitOutput, apperrors.ApplicationError) {
	app := u.factory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	// Build exercise map for quick lookup
	exerciseMap := map[string]domain.Exercise{}
	for _, pse := range ps.Exercises {
		exerciseMap[pse.Exercise.ID] = pse.Exercise
	}

	profile, _ := app.Repositories.UserProfile.Get(ctx, studentID)
	assistantCfg := assistant.Config{}
	if profile != nil {
		assistantCfg.BaseURL = profile.AssistantBaseURL
		assistantCfg.APIKey = profile.AssistantAPIKey
	}

	correct := 0
	total := len(input.Attempts)
	totalHints := 0
	totalTime := 0

	for _, attempt := range input.Attempts {
		ex, ok := exerciseMap[attempt.ExerciseID]
		isCorrect := false
		answerText := attempt.AnswerText

		if attempt.CanvasData != "" && app.AssistantService != nil && app.AssistantService.IsConfigured(assistantCfg) {
			if recognizedText, recognizeErr := app.AssistantService.AnalyzeCanvas(ctx, assistantCfg, attempt.CanvasData, ex.CorrectAnswer); recognizeErr == nil {
				normalizedRecognized := normalizeCanvasAnswer(recognizedText)
				if normalizedRecognized != "" && normalizedRecognized != "UNREADABLE" {
					answerText = normalizedRecognized
				}
			}
		}

		if ok && ex.CorrectAnswer != "" {
			isCorrect = strings.EqualFold(
				normalizeCanvasAnswer(answerText),
				normalizeCanvasAnswer(ex.CorrectAnswer),
			)
		}

		score := 0.0
		if isCorrect {
			correct++
			score = 100.0
		}

		totalHints += attempt.HintsUsed
		totalTime += attempt.TimeSpentSeconds

		attemptID, _ := app.Repositories.StudentAttempt.Create(ctx, domain.StudentAttempt{
			StudentID:       studentID,
			ExerciseID:      attempt.ExerciseID,
			PracticeSheetID: sheetID,
			AnswerText:      answerText,
			IsCorrect:       isCorrect,
			Score:           score,
			TimeSpentSecs:   attempt.TimeSpentSeconds,
			HintsUsed:       attempt.HintsUsed,
		})

		if attempt.CanvasData != "" && attemptID != "" {
			app.Repositories.StudentAttempt.SaveCanvasWork(ctx, attemptID, attempt.CanvasData)
		}
	}

	sheetScore := 0.0
	if total > 0 {
		sheetScore = float64(correct) / float64(total) * 100
	}

	kumon := app.KumonStrategy
	currentProgress, _ := app.Repositories.StudentProgress.Get(ctx, studentID, ps.TopicID)
	currentScore := 0.0
	currentLevel := 1
	prevTotal := 0
	prevCorrect := 0
	if currentProgress != nil {
		currentScore = currentProgress.MasteryScore
		currentLevel = currentProgress.CurrentLevel
		prevTotal = currentProgress.TotalAttempts
		prevCorrect = currentProgress.CorrectAttempts
	}

	var newMastery float64
	var shouldLevelUp bool
	var shouldRepeat bool
	var nextLevel int
	var recommendation string

	const levelTestPassThreshold = 75.0

	if ps.SheetType == "level_test" {
		// Level test: pass/fail by fixed threshold, force level up on pass
		shouldLevelUp = sheetScore >= levelTestPassThreshold
		shouldRepeat = !shouldLevelUp
		if shouldLevelUp {
			nextLevel = currentLevel + 1
			newMastery = sheetScore
			recommendation = "¡Aprobaste la prueba! Nivel " + itoa(nextLevel) + " desbloqueado."
		} else {
			nextLevel = currentLevel
			newMastery = currentScore
			recommendation = "Necesitás al menos 75% para pasar de nivel. ¡Seguí practicando!"
		}
	} else {
		newMastery = kumon.CalculateMasteryScore(strategy.MasteryInput{
			TotalAttempts:    total,
			CorrectAttempts:  correct,
			HintsUsed:        totalHints,
			TimeSpentSeconds: totalTime,
			CurrentScore:     currentScore,
		})
		rec := kumon.GenerateNextPracticeRecommendation(newMastery, currentLevel)
		shouldLevelUp = kumon.ShouldLevelUp(newMastery)
		shouldRepeat = kumon.ShouldRepeatTopic(newMastery)
		recommendation = rec.Message
		nextLevel = currentLevel
		if shouldLevelUp {
			nextLevel = currentLevel + 1
		}
	}

	app.Repositories.StudentProgress.Upsert(ctx, domain.StudentTopicProgress{
		StudentID:       studentID,
		TopicID:         ps.TopicID,
		StrategyID:      ps.StrategyID,
		MasteryScore:    newMastery,
		CurrentLevel:    nextLevel,
		TotalAttempts:   prevTotal + total,
		CorrectAttempts: prevCorrect + correct,
	})

	// On level test pass, also advance the student's course-level progress
	if ps.SheetType == "level_test" && shouldLevelUp {
		app.Repositories.CourseProgress.Upsert(ctx, studentID, ps.CourseID, nextLevel)
	}

	return &SubmitOutput{Data: SubmitResult{
		Score:          sheetScore,
		Correct:        correct,
		Total:          total,
		MasteryScore:   newMastery,
		Recommendation: recommendation,
		ShouldLevelUp:  shouldLevelUp,
		ShouldRepeat:   shouldRepeat,
		NextLevel:      nextLevel,
	}}, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}

func normalizeCanvasAnswer(value string) string {
	normalized := strings.TrimSpace(value)
	normalized = strings.ReplaceAll(normalized, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\t", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")
	return normalized
}
