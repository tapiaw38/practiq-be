package practicesheet

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/assistant"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	SubmitUsecase interface {
		Execute(context.Context, string, string, SubmitInput) (*SubmitOutput, apperrors.ApplicationError)
	}

	submitUsecase struct {
		contextFactory appcontext.Factory
	}

	AttemptInput struct {
		ExerciseID       string `json:"exercise_id"`
		AnswerText       string `json:"answer_text"`
		CanvasData       string `json:"canvas_data"` // base64 PNG for canvas/handwritten exercises
		TimeSpentSeconds int    `json:"time_spent_seconds"`
		HintsUsed        int    `json:"hints_used"`
	}

	SubmitInput struct {
		Attempts []AttemptInput `json:"attempts"`
	}

	SubmitOutput struct {
		Data SubmitResult `json:"data"`
	}
)

func NewSubmitUsecase(contextFactory appcontext.Factory) SubmitUsecase {
	return &submitUsecase{contextFactory: contextFactory}
}

func (u *submitUsecase) Execute(ctx context.Context, sheetID, studentID string, input SubmitInput) (*SubmitOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	exerciseMap := map[string]domain.Exercise{}
	derivedTopicID := ps.TopicID
	for _, pse := range ps.Exercises {
		exerciseMap[pse.Exercise.ID] = pse.Exercise
		if derivedTopicID == "" && pse.Exercise.TopicID != "" {
			derivedTopicID = pse.Exercise.TopicID
		}
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
	resultAIFeedback := ""

	for _, attempt := range input.Attempts {
		ex, ok := exerciseMap[attempt.ExerciseID]
		isCorrect := false
		answerText := attempt.AnswerText
		imageURL := ""
		aiFeedback := ""
		hasTextAnswer := strings.TrimSpace(answerText) != ""
		hasCanvasAnswer := strings.TrimSpace(attempt.CanvasData) != ""

		if hasCanvasAnswer && app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
			if recognizedText, recognizeErr := app.Integrations.AssistantGateway.AnalyzeCanvas(ctx, assistantCfg, attempt.CanvasData, ex.CorrectAnswer); recognizeErr == nil {
				normalizedRecognized := normalizeCanvasAnswer(recognizedText)
				if normalizedRecognized != "" && normalizedRecognized != "UNREADABLE" {
					answerText = normalizedRecognized
					hasTextAnswer = true
				}
			}
		}

		if ok {
			normalizedCorrect := normalizeCanvasAnswer(ex.CorrectAnswer)
			isCorrect = strings.EqualFold(
				normalizeCanvasAnswer(answerText),
				normalizedCorrect,
			)

			canEvaluateWithAI := app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) && hasTextAnswer && !isDataURIAnswer(answerText)
			if canEvaluateWithAI {
				if evaluation, aiErr := app.Integrations.AssistantGateway.EvaluatePracticeAnswer(ctx, assistantCfg, ex.Question, ex.CorrectAnswer, answerText); aiErr == nil {
					isCorrect = evaluation.IsCorrect
					aiFeedback = evaluation.Feedback
					if resultAIFeedback == "" && aiFeedback != "" {
						resultAIFeedback = aiFeedback
					}
				}
			}
		}

		if hasCanvasAnswer && app.ImageStorage != nil {
			canvasData := normalizeCanvasDataURI(attempt.CanvasData)
			if uploaded, uploadErr := app.ImageStorage.UploadDataURI(ctx, "practice", studentID, canvasData); uploadErr == nil {
				imageURL = uploaded
			} else {
				log.Printf("[image_storage] practice attempt upload failed student_id=%s exercise_id=%s err=%v", studentID, attempt.ExerciseID, uploadErr)
			}
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
			ImageURL:        imageURL,
			AIFeedback:      aiFeedback,
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

	kumon := domain.NewKumonStrategy()
	currentProgress, _ := app.Repositories.StudentProgress.Get(ctx, studentID, derivedTopicID)
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
		shouldLevelUp = sheetScore >= levelTestPassThreshold
		shouldRepeat = !shouldLevelUp
		if shouldLevelUp {
			nextLevel = currentLevel + 1
			newMastery = sheetScore
			recommendation = "¡Aprobaste la prueba! Nivel " + strconv.Itoa(nextLevel) + " desbloqueado."
		} else {
			nextLevel = currentLevel
			newMastery = currentScore
			recommendation = "Necesitás al menos 75% para pasar de nivel. ¡Seguí practicando!"
		}
	} else {
		newMastery = kumon.CalculateMasteryScore(domain.MasteryInput{
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

	newStreak := calcStreak(currentProgress)
	now := time.Now()

	if derivedTopicID != "" {
		if err := app.Repositories.StudentProgress.Upsert(ctx, domain.StudentTopicProgress{
			StudentID:       studentID,
			TopicID:         derivedTopicID,
			StrategyID:      ps.StrategyID,
			MasteryScore:    newMastery,
			CurrentLevel:    nextLevel,
			TotalAttempts:   prevTotal + total,
			CorrectAttempts: prevCorrect + correct,
			StreakDays:      newStreak,
			LastPracticedAt: &now,
		}); err != nil {
			// Log error but don't fail the submit - progress update failure should not break submit
			_ = err
		}
	}

	if ps.SheetType == "level_test" && shouldLevelUp {
		app.Repositories.CourseProgress.Upsert(ctx, studentID, ps.CourseID, nextLevel)
	}

	return &SubmitOutput{Data: toSubmitOutputData(sheetScore, correct, total, newMastery, recommendation, resultAIFeedback, shouldLevelUp, shouldRepeat, nextLevel)}, nil
}

func normalizeCanvasAnswer(value string) string {
	normalized := strings.TrimSpace(value)
	normalized = strings.ReplaceAll(normalized, "\n", " ")
	normalized = strings.ReplaceAll(normalized, "\t", " ")
	normalized = strings.Join(strings.Fields(normalized), " ")
	return normalized
}

func isDataURIAnswer(value string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	return strings.HasPrefix(trimmed, "data:image/")
}

func normalizeCanvasDataURI(value string) string {
	trimmed := strings.TrimSpace(value)
	if strings.HasPrefix(strings.ToLower(trimmed), "data:image/") {
		return trimmed
	}
	return "data:image/png;base64," + trimmed
}

func calcStreak(current *domain.StudentTopicProgress) int {
	if current == nil || current.LastPracticedAt == nil {
		return 1
	}
	now := time.Now().UTC()
	last := current.LastPracticedAt.UTC()

	nowDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	lastDay := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, time.UTC)

	diff := int(nowDay.Sub(lastDay).Hours() / 24)

	switch {
	case diff == 0:
		return current.StreakDays
	case diff == 1:
		return current.StreakDays + 1
	default:
		return 1
	}
}
