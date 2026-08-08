package practicesheet

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
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
		// Attachment answers: URL returned by POST /uploads plus its metadata.
		AttachmentURL         string `json:"attachment_url"`
		AttachmentName        string `json:"attachment_name"`
		AttachmentContentType string `json:"attachment_content_type"`
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

	hasAccess, err := studentHasCourseAccess(ctx, app, studentID, ps.CourseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if !hasAccess {
		return nil, apperrors.NewForbiddenError()
	}
	if appErr := ensureSheetIsOpen(ctx, app, ps, studentID, false); appErr != nil {
		return nil, appErr
	}

	// Get course for grade context
	course, _ := app.Repositories.Course.Get(ctx, ps.CourseID)
	gradeName := ""
	if course != nil {
		gradeName = course.GradeName
	}

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
	resultAIFeedback := ""
	exerciseResults := make([]ExerciseResultData, 0, total)

	// Track progress per topic
	topicStats := make(map[string]struct{ correct, total int })

	for _, attempt := range input.Attempts {
		ex, ok := exerciseMap[attempt.ExerciseID]
		isCorrect := false
		answerText := attempt.AnswerText
		imageURL := ""
		aiFeedback := ""
		hasTextAnswer := strings.TrimSpace(answerText) != ""
		hasCanvasAnswer := strings.TrimSpace(attempt.CanvasData) != ""

		if hasCanvasAnswer && app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
			normalizedCanvas := normalizeCanvasDataURI(attempt.CanvasData)
			if recognizedText, recognizeErr := app.Integrations.AssistantGateway.AnalyzeCanvas(ctx, assistantCfg, normalizedCanvas, ex.CorrectAnswer); recognizeErr == nil {
				normalizedRecognized := normalizeCanvasAnswer(recognizedText)
				if normalizedRecognized != "" && normalizedRecognized != "UNREADABLE" {
					answerText = normalizedRecognized
					hasTextAnswer = true
				}
			}
		}

		hasAttachment := strings.TrimSpace(attempt.AttachmentURL) != ""
		needsTeacherReview := false

		if ok && ex.Type == exerciseTypeAttachment {
			if !hasAttachment {
				// Nothing was uploaded: that is simply an unanswered exercise.
				needsTeacherReview = false
			} else {
				outcome := evaluateAttachment(ctx, app, assistantCfg, ex, gradeName,
					attempt.AttachmentURL, attempt.AttachmentName, attempt.AttachmentContentType)
				isCorrect = outcome.IsCorrect
				aiFeedback = outcome.Feedback
				needsTeacherReview = outcome.NeedsReview
				if resultAIFeedback == "" && aiFeedback != "" {
					resultAIFeedback = aiFeedback
				}
			}
		} else if ok {
			normalizedCorrect := normalizeCanvasAnswer(ex.CorrectAnswer)
			isCorrect = strings.EqualFold(
				normalizeCanvasAnswer(answerText),
				normalizedCorrect,
			)

			canEvaluateWithAI := app.Integrations.AssistantGateway != nil && app.Integrations.AssistantGateway.IsConfigured(assistantCfg) && hasTextAnswer && !isDataURIAnswer(answerText)
			if canEvaluateWithAI {
				if evaluation, aiErr := app.Integrations.AssistantGateway.EvaluatePracticeAnswer(ctx, assistantCfg, ex.Question, ex.CorrectAnswer, answerText, gradeName); aiErr == nil {
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
		switch {
		case needsTeacherReview:
			// Not graded yet: drop it from the denominator so an unread file
			// cannot fail the student on its own.
			total--
		case isCorrect:
			correct++
			score = 100.0
		}

		// Track per-topic stats
		if ok && ex.TopicID != "" && !needsTeacherReview {
			stats := topicStats[ex.TopicID]
			stats.total++
			if isCorrect {
				stats.correct++
			}
			topicStats[ex.TopicID] = stats
		}

		totalHints += attempt.HintsUsed
		totalTime += attempt.TimeSpentSeconds

		attemptID, _ := app.Repositories.StudentAttempt.Create(ctx, domain.StudentAttempt{
			StudentID:             studentID,
			ExerciseID:            attempt.ExerciseID,
			PracticeSheetID:       sheetID,
			AnswerText:            answerText,
			ImageURL:              imageURL,
			AIFeedback:            aiFeedback,
			IsCorrect:             isCorrect,
			Score:                 score,
			TimeSpentSecs:         attempt.TimeSpentSeconds,
			HintsUsed:             attempt.HintsUsed,
			AttachmentURL:         attempt.AttachmentURL,
			AttachmentName:        attempt.AttachmentName,
			AttachmentContentType: attempt.AttachmentContentType,
			NeedsTeacherReview:    needsTeacherReview,
		})

		if attempt.CanvasData != "" && attemptID != "" {
			app.Repositories.StudentAttempt.SaveCanvasWork(ctx, attemptID, attempt.CanvasData)
		}

		// Collect per-exercise result
		correctAnswer := ""
		if ok {
			correctAnswer = ex.CorrectAnswer
		}
		exerciseResults = append(exerciseResults, ExerciseResultData{
			ExerciseID:         attempt.ExerciseID,
			IsCorrect:          isCorrect,
			StudentAnswer:      answerText,
			CorrectAnswer:      correctAnswer,
			AIFeedback:         aiFeedback,
			NeedsTeacherReview: needsTeacherReview,
		})
	}

	sheetScore := 0.0
	if total > 0 {
		sheetScore = float64(correct) / float64(total) * 100
	}
	// Every answer is awaiting review: there is no score to act on yet.
	allPendingReview := total <= 0 && len(input.Attempts) > 0

	kumon := domain.NewKumonStrategy()

	// For level tests or single-topic sheets, use overall score for level progression
	derivedTopicID := ps.TopicID
	if derivedTopicID == "" && len(topicStats) == 1 {
		for topicID := range topicStats {
			derivedTopicID = topicID
		}
	}

	currentProgress, _ := app.Repositories.StudentProgress.Get(ctx, studentID, derivedTopicID)
	currentScore := 0.0
	currentLevel := 1
	if currentProgress != nil {
		currentScore = currentProgress.MasteryScore
		currentLevel = currentProgress.CurrentLevel
	}

	var newMastery float64
	var shouldLevelUp bool
	var shouldRepeat bool
	var nextLevel int
	var recommendation string

	const levelTestPassThreshold = 75.0

	switch {
	case allPendingReview:
		// Nothing gradeable came back yet, so level progression waits for the
		// teacher rather than acting on a score of zero.
		shouldLevelUp = false
		shouldRepeat = false
		nextLevel = currentLevel
		newMastery = currentScore
		recommendation = "Tu entrega quedó pendiente de la revisión del docente."
	case ps.SheetType == "level_test":
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
	default:
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

	now := time.Now()

	// Write progress for each topic
	for topicID, stats := range topicStats {
		topicProgress, _ := app.Repositories.StudentProgress.Get(ctx, studentID, topicID)
		topicCurrentScore := 0.0
		topicPrevTotal := 0
		topicPrevCorrect := 0
		if topicProgress != nil {
			topicCurrentScore = topicProgress.MasteryScore
			topicPrevTotal = topicProgress.TotalAttempts
			topicPrevCorrect = topicProgress.CorrectAttempts
		}

		topicMastery := kumon.CalculateMasteryScore(domain.MasteryInput{
			TotalAttempts:    stats.total,
			CorrectAttempts:  stats.correct,
			HintsUsed:        0,
			TimeSpentSeconds: 0,
			CurrentScore:     topicCurrentScore,
		})

		topicStreak := calcStreak(topicProgress)

		levelToSave := 1
		if topicProgress != nil {
			levelToSave = topicProgress.CurrentLevel
		}
		if topicID == derivedTopicID {
			if shouldLevelUp {
				levelToSave = nextLevel
			} else {
				levelToSave = currentLevel
			}
		}

		if err := app.Repositories.StudentProgress.Upsert(ctx, domain.StudentTopicProgress{
			StudentID:       studentID,
			TopicID:         topicID,
			StrategyID:      ps.StrategyID,
			MasteryScore:    topicMastery,
			CurrentLevel:    levelToSave,
			TotalAttempts:   topicPrevTotal + stats.total,
			CorrectAttempts: topicPrevCorrect + stats.correct,
			StreakDays:      topicStreak,
			LastPracticedAt: &now,
		}); err != nil {
			// Log error but don't fail the submit
			_ = err
		}
	}

	if ps.SheetType == "level_test" && shouldLevelUp {
		app.Repositories.CourseProgress.Upsert(ctx, studentID, ps.CourseID, nextLevel)
	}

	result := toSubmitOutputData(sheetScore, correct, total, newMastery, recommendation, resultAIFeedback, shouldLevelUp, shouldRepeat, nextLevel, exerciseResults)
	result.PendingReview = allPendingReview
	return &SubmitOutput{Data: result}, nil
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
