package practicesheet

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
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
		CanvasData       string `json:"canvas_data"`
		TimeSpentSeconds int    `json:"time_spent_seconds"`
		HintsUsed        int    `json:"hints_used"`

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

	if appErr := school.EnsureCourseAcceptsWork(ctx, app, ps.CourseID); appErr != nil {
		return nil, appErr
	}

	if appErr := school.EnsureStudentCanWork(ctx, app, studentID, ps.CourseID); appErr != nil {
		return nil, appErr
	}
	if appErr := ensureSheetIsOpen(ctx, app, ps, studentID, false); appErr != nil {
		return nil, appErr
	}
	if ps.SheetType == sheetTypeLevelTest {
		if appErr := validateLevelTestAttempts(ps.Exercises, input.Attempts); appErr != nil {
			return nil, appErr
		}
		if appErr := ensureWithinTimeLimit(ctx, app, *ps, studentID); appErr != nil {
			return nil, appErr
		}
		allowed := ps.AttemptsAllowed()
		claimed, claimErr := app.Repositories.StudentAttempt.ClaimLevelTestSubmission(ctx, studentID, sheetID, allowed)
		if claimErr != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, claimErr)
		}
		if !claimed {
			if allowed == 1 {
				return nil, apperrors.NewBadRequestError("this level test was already submitted")
			}
			return nil, apperrors.NewBadRequestError(fmt.Sprintf("no attempts left: this level test allows %d", allowed))
		}
	}

	course, _ := app.Repositories.Course.Get(ctx, ps.CourseID)
	gradeName := ""
	if course != nil {
		gradeName = course.GradeName
	}

	exerciseMap := map[string]domain.Exercise{}
	for _, pse := range ps.Exercises {
		exerciseMap[pse.Exercise.ID] = pse.Exercise
	}

	assistantCfg := assistantcfg.Resolve(ctx, app)

	teacherGrades := teacherGradesSheet(ps.SheetType)

	correct := 0
	total := len(input.Attempts)

	hasPendingReview := false
	totalHints := 0
	totalTime := 0

	var persistenceErr error
	resultAIFeedback := ""
	exerciseResults := make([]ExerciseResultData, 0, total)
	xpExercises := make([]xpExercise, 0, total)

	studentLoc := domain.StudentLocation("")
	if profile, profileErr := app.Repositories.UserProfile.Get(ctx, studentID); profileErr == nil && profile != nil {
		studentLoc = domain.StudentLocation(profile.Timezone)
	}

	topicStats := make(map[string]struct{ correct, total int })

	for _, attempt := range input.Attempts {
		ex, ok := exerciseMap[attempt.ExerciseID]
		isCorrect := false
		answerText := attempt.AnswerText
		imageURL := ""
		aiFeedback := ""
		hasTextAnswer := strings.TrimSpace(answerText) != ""
		hasCanvasAnswer := strings.TrimSpace(attempt.CanvasData) != ""

		if strings.TrimSpace(attempt.AttachmentURL) != "" &&
			(app.ImageStorage == nil ||
				!app.ImageStorage.OwnsFileURL(attempt.AttachmentURL, attachmentsFolder, studentID)) {
			log.Printf("[practice_attachment] rejected foreign attachment student_id=%s url=%q", studentID, attempt.AttachmentURL)
			attempt.AttachmentURL = ""
			attempt.AttachmentName = ""
			attempt.AttachmentContentType = ""
		}

		hasAttachment := strings.TrimSpace(attempt.AttachmentURL) != ""

		hasStatementMedia := ok && ex.MediaURL() != ""
		statementMediaNeedsReview := needsReviewForStatementMedia(ex, hasTextAnswer, hasCanvasAnswer, hasAttachment)
		canvasUnreadable := false

		canvasAwaitsOCR := hasCanvasAnswer && !hasStatementMedia
		assistantReady := app.Integrations.AssistantGateway != nil &&
			app.Integrations.AssistantGateway.IsConfigured(assistantCfg)

		if canvasAwaitsOCR && assistantReady {
			normalizedCanvas := normalizeCanvasDataURI(attempt.CanvasData)
			if recognizedText, recognizeErr := app.Integrations.AssistantGateway.AnalyzeCanvas(ctx, assistantCfg, normalizedCanvas, ex.CorrectAnswer); recognizeErr == nil {
				normalizedRecognized := normalizeCanvasAnswer(recognizedText)
				if normalizedRecognized != "" && normalizedRecognized != "UNREADABLE" {
					answerText = normalizedRecognized
					hasTextAnswer = true
				} else if !hasTextAnswer {
					canvasUnreadable = true
				}
			} else if !hasTextAnswer {

				canvasUnreadable = true
			}
		}

		if transcriptionUnavailable(canvasAwaitsOCR, hasTextAnswer, assistantReady) {
			canvasUnreadable = true
		}

		ungraded := canvasUnreadable || statementMediaNeedsReview
		if canvasUnreadable {
			answerText = "UNREADABLE"
			aiFeedback = unreadableCanvasFeedback(teacherGrades)
		} else if statementMediaNeedsReview {
			aiFeedback = statementMediaFeedback(teacherGrades)
		}
		var aiSuggestion *bool

		if statementMediaNeedsReview {

		} else if ok && ex.Type == exerciseTypeAttachment {
			if !hasAttachment {

				ungraded = false
			} else {
				outcome := evaluateAttachment(ctx, app, assistantCfg, ex, gradeName,
					attempt.AttachmentURL, attempt.AttachmentName, teacherGrades)
				isCorrect = outcome.IsCorrect
				aiSuggestion = outcome.AISuggestedCorrect
				aiFeedback = outcome.Feedback
				ungraded = outcome.Ungraded
				if resultAIFeedback == "" && aiFeedback != "" {
					resultAIFeedback = aiFeedback
				}
			}
		} else if ok && ex.Type == exerciseTypeFillBlanks {

			isCorrect = blanksAnswersMatch(answerText, ex.CorrectAnswer)
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
				} else if isCorrect, ungraded = unresolvedEvaluation(isCorrect); ungraded {
					log.Printf("[practice_submit] evaluation unavailable, leaving ungraded student_id=%s exercise_id=%s err=%v", studentID, attempt.ExerciseID, aiErr)
					if aiFeedback == "" {
						aiFeedback = evaluationUnavailableFeedback(teacherGrades)
					}
					if resultAIFeedback == "" {
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

				imageURL = canvasData
			}
		}

		needsTeacherReview := ungraded && teacherGrades

		score := 0.0
		switch {
		case ungraded:

			total--
			if needsTeacherReview {
				hasPendingReview = true
			}
		case isCorrect:
			correct++
			score = 100.0
		}

		if ok && ex.TopicID != "" && !ungraded {
			stats := topicStats[ex.TopicID]
			stats.total++
			if isCorrect {
				stats.correct++
			}
			topicStats[ex.TopicID] = stats
		}

		totalHints += attempt.HintsUsed
		totalTime += attempt.TimeSpentSeconds

		attemptID, createErr := app.Repositories.StudentAttempt.Create(ctx, domain.StudentAttempt{
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
			NotGraded:             ungraded,
			AIIsCorrect:           aiSuggestion,
		})

		if createErr != nil {
			log.Printf("[practice_submit] could not persist attempt student_id=%s sheet_id=%s err=%v", studentID, sheetID, createErr)
			persistenceErr = createErr
			break
		}

		if attempt.CanvasData != "" && attemptID != "" {
			if saveErr := app.Repositories.StudentAttempt.SaveCanvasWork(ctx, attemptID, attempt.CanvasData); saveErr != nil {
				log.Printf("[practice_submit] could not persist canvas student_id=%s sheet_id=%s attempt_id=%s err=%v", studentID, sheetID, attemptID, saveErr)
				persistenceErr = saveErr
				break
			}
		}

		if ok {
			xpExercises = append(xpExercises, xpExercise{ExerciseID: attempt.ExerciseID, Correct: isCorrect, Ungraded: ungraded})
		}

		if ps.SheetType == sheetTypeLevelTest {
			continue
		}
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
			NotGraded:          ungraded,
		})
	}

	if persistenceErr != nil {
		if ps.SheetType != sheetTypeLevelTest {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetSubmitError, persistenceErr)
		}

		cleanupCtx, cancelCleanup := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
		defer cancelCleanup()

		if cleanupErr := app.Repositories.StudentAttempt.DeleteBySheet(cleanupCtx, studentID, sheetID); cleanupErr != nil {
			log.Printf("[practice_submit] could not remove partial level test student_id=%s sheet_id=%s err=%v", studentID, sheetID, cleanupErr)
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetSubmitError, cleanupErr)
		}
		if releaseErr := app.Repositories.StudentAttempt.ReleaseLevelTestSubmission(cleanupCtx, studentID, sheetID); releaseErr != nil {
			log.Printf("[practice_submit] could not release the level test claim student_id=%s sheet_id=%s err=%v", studentID, sheetID, releaseErr)
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetSubmitError, releaseErr)
		}
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetSubmitError, persistenceErr)
	}

	sheetScore := 0.0
	if total > 0 {
		sheetScore = float64(correct) / float64(total) * 100
	}

	allUngraded := total <= 0 && len(input.Attempts) > 0

	kumon := domain.NewKumonStrategy()

	derivedTopicID := ps.TopicID
	if derivedTopicID == "" && len(topicStats) == 1 {
		for topicID := range topicStats {
			derivedTopicID = topicID
		}
	}

	currentScore := 0.0
	currentLevel := 1
	if ps.SheetType == "level_test" {

		currentLevel = ps.Level
	} else {
		currentProgress, _ := app.Repositories.StudentProgress.Get(ctx, studentID, derivedTopicID)
		if currentProgress != nil {
			currentScore = currentProgress.MasteryScore
			currentLevel = currentProgress.CurrentLevel
		}
	}

	var newMastery float64
	var shouldLevelUp bool
	var shouldRepeat bool
	var nextLevel int
	var recommendation string

	const levelTestPassThreshold = 75.0

	switch {
	case allUngraded:

		shouldLevelUp = false
		shouldRepeat = false
		nextLevel = currentLevel
		newMastery = currentScore
		if hasPendingReview {
			recommendation = "Tu entrega quedó pendiente de la revisión del docente."
		} else {
			recommendation = "No pudimos corregir tus respuestas automáticamente, así que tu progreso quedó como estaba."
		}
	case ps.SheetType == sheetTypeLevelTest && hasPendingReview:

		shouldLevelUp = false
		shouldRepeat = false
		nextLevel = currentLevel
		newMastery = currentScore
		recommendation = "Tu prueba quedó esperando la corrección del docente."
	case ps.SheetType == sheetTypeLevelTest:
		shouldLevelUp = sheetScore >= levelTestPassThreshold
		shouldRepeat = !shouldLevelUp
		if shouldLevelUp {
			nextLevel = ps.Level + 1
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

		topicStreak := calcStreak(topicProgress, studentLoc)

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

			_ = err
		}
	}

	if ps.SheetType == "level_test" && shouldLevelUp {
		if err := app.Repositories.CourseProgress.Upsert(ctx, studentID, ps.CourseID, nextLevel); err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetSubmitError, err)
		}
	}

	if ps.SheetType == sheetTypeLevelTest {
		resultAIFeedback = ""
	}
	result := toSubmitOutputData(sheetScore, correct, total, newMastery, recommendation, resultAIFeedback, shouldLevelUp, shouldRepeat, nextLevel, exerciseResults)
	xp := awardPracticeXP(ctx, app, studentID, ps, input.Attempts, xpExercises, ps.SheetType == sheetTypeLevelTest && shouldLevelUp && !hasPendingReview)
	result.XPGained, result.CourseXP, result.XPBreakdown = xp.Gained, xp.Balance, xp.Breakdown
	if progress, err := app.Repositories.StudentProgress.ListByStudent(ctx, studentID); err != nil {
		log.Printf("[practice_sheet] could not refresh streak student_id=%s err=%v", studentID, err)
	} else {
		result.StreakDays = domain.CurrentStreak(progress, studentLoc)
	}
	result.PendingReview = hasPendingReview
	return &SubmitOutput{Data: result}, nil
}

func validateLevelTestAttempts(exercises []domain.PracticeSheetExercise, attempts []AttemptInput) apperrors.ApplicationError {
	if len(attempts) != len(exercises) {
		return apperrors.NewBadRequestError("a level test must include every exercise exactly once")
	}

	exerciseIDs := make(map[string]struct{}, len(exercises))
	for _, pse := range exercises {
		exerciseIDs[pse.Exercise.ID] = struct{}{}
	}
	seen := make(map[string]struct{}, len(attempts))
	for _, attempt := range attempts {
		if _, exists := exerciseIDs[attempt.ExerciseID]; !exists {
			return apperrors.NewBadRequestError("the level test contains an unknown exercise")
		}
		if _, duplicate := seen[attempt.ExerciseID]; duplicate {
			return apperrors.NewBadRequestError("each level test exercise can be submitted only once")
		}
		seen[attempt.ExerciseID] = struct{}{}
	}
	return nil
}

func needsReviewForStatementMedia(ex domain.Exercise, hasTextAnswer, hasCanvasAnswer, hasAttachment bool) bool {
	return ex.MediaURL() != "" && (hasTextAnswer || hasCanvasAnswer || hasAttachment)
}

func transcriptionUnavailable(canvasAwaitsOCR, hasTextAnswer, assistantReady bool) bool {
	return canvasAwaitsOCR && !hasTextAnswer && !assistantReady
}

func teacherGradesSheet(sheetType string) bool {
	return sheetType == sheetTypeLevelTest
}

func unresolvedEvaluation(textMatches bool) (isCorrect, ungraded bool) {
	if textMatches {
		return true, false
	}
	return false, true
}

func evaluationUnavailableFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "No pudimos evaluar tu respuesta en este momento, así que la va a revisar el docente."
	}
	return "No pudimos evaluar tu respuesta en este momento, así que no cuenta en tu puntaje."
}

func unreadableCanvasFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "No pudimos leer tu respuesta escrita, así que la va a corregir el docente."
	}
	return "No pudimos leer tu respuesta escrita, así que no cuenta en tu puntaje. Intentá escribirla más clara la próxima vez."
}

func statementMediaFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "Tu respuesta quedó pendiente de la revisión del docente porque este ejercicio incluye material visual o de audio."
	}
	return "No pudimos corregir esta respuesta automáticamente porque el ejercicio incluye material visual o de audio, así que no cuenta en tu puntaje."
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

func calcStreak(current *domain.StudentTopicProgress, loc *time.Location) int {
	if current == nil || current.LastPracticedAt == nil {
		return 1
	}

	switch domain.DaysBetween(*current.LastPracticedAt, time.Now(), loc) {
	case 0:
		return current.StreakDays
	case 1:
		return current.StreakDays + 1
	default:
		return 1
	}
}

func ensureWithinTimeLimit(ctx context.Context, app *appcontext.Context, ps domain.PracticeSheet, studentID string) apperrors.ApplicationError {
	if ps.TimeLimitMinutes == nil {
		return nil
	}
	_, startedAt, err := app.Repositories.StudentAttempt.LevelTestProgress(ctx, studentID, ps.ID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	deadline := ps.Deadline(startedAt)
	if deadline == nil || !submissionReferenceTime(ctx).After(*deadline) {
		return nil
	}
	if err := app.Repositories.StudentAttempt.CloseExpiredLevelTest(ctx, studentID, ps.ID, *deadline); err != nil {
		return apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	return apperrors.NewBadRequestError(fmt.Sprintf("time is up: this level test allows %d minutes", *ps.TimeLimitMinutes))
}
