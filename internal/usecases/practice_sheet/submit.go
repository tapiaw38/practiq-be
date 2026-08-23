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
	if ps.SheetType == sheetTypeLevelTest {
		if appErr := validateLevelTestAttempts(ps.Exercises, input.Attempts); appErr != nil {
			return nil, appErr
		}
		claimed, claimErr := app.Repositories.StudentAttempt.ClaimLevelTestSubmission(ctx, studentID, sheetID)
		if claimErr != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, claimErr)
		}
		if !claimed {
			return nil, apperrors.NewBadRequestError("this level test was already submitted")
		}
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
	// The assistant reads this answer on the course's behalf, not the student's.
	// A student with no credentials of their own left every handwritten answer
	// untranscribed, which on a level test cost them the exercise. The notebook
	// flow closed the same hole by resolving the teacher's config.
	if course != nil && app.Integrations.AssistantGateway != nil &&
		!app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
		assistantCfg = teacherAssistantConfig(ctx, app, course.TeacherID)
	}

	// Only a level test is corrected by a teacher (the homework notebook has its
	// own review flow). A practice is graded on the spot with whatever the
	// assistant could read: it never enters the teacher's queue and never waits
	// on one.
	teacherGrades := teacherGradesSheet(ps.SheetType)

	correct := 0
	total := len(input.Attempts)
	// Any answer waiting for a teacher blocks promotion on a level test.
	hasPendingReview := false
	totalHints := 0
	totalTime := 0
	// A level test claims its single submission before any of this is written.
	// Any persistence failure must compensate every saved attempt and release
	// the claim, otherwise a partial delivery could be scored or lock the user.
	var persistenceErr error
	resultAIFeedback := ""
	exerciseResults := make([]ExerciseResultData, 0, total)

	// The streak counts calendar days in the student's own zone, so it has to
	// be resolved before any topic is updated.
	studentLoc := domain.StudentLocation("")
	if profile, profileErr := app.Repositories.UserProfile.Get(ctx, studentID); profileErr == nil && profile != nil {
		studentLoc = domain.StudentLocation(profile.Timezone)
	}

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
		// Only a file this student uploaded counts. Without the ownership check
		// a crafted attachment_url pointed at any object in the bucket — another
		// student's delivery or a course material — and the submit flow both
		// forwarded its contents to the assistant and stored it on this attempt,
		// where the teacher's review view later presigns it.
		if strings.TrimSpace(attempt.AttachmentURL) != "" &&
			(app.ImageStorage == nil ||
				!app.ImageStorage.OwnsFileURL(attempt.AttachmentURL, attachmentsFolder, studentID)) {
			log.Printf("[practice_attachment] rejected foreign attachment student_id=%s url=%q", studentID, attempt.AttachmentURL)
			attempt.AttachmentURL = ""
			attempt.AttachmentName = ""
			attempt.AttachmentContentType = ""
		}

		hasAttachment := strings.TrimSpace(attempt.AttachmentURL) != ""
		// Gillie does not participate in grading. A statement image or audio can
		// change what the answer means, so do not grade it from the text-only
		// question and answer fields. Keep a submitted answer for teacher review.
		hasStatementMedia := ok && ex.MediaURL() != ""
		statementMediaNeedsReview := needsReviewForStatementMedia(ex, hasTextAnswer, hasCanvasAnswer, hasAttachment)
		canvasUnreadable := false

		// Handwriting is only an image until the assistant reads it, so this is
		// the one path that can turn it into something comparable.
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
				// Do not turn an OCR failure into a fabricated wrong answer.
				canvasUnreadable = true
			}
		}
		// Nothing ever read the handwriting: no assistant at all, or none
		// configured for the student or their teacher. Without this the empty
		// answerText fell through to the comparison below and scored as wrong —
		// the student lost the exercise for a transcription that never ran, and
		// on a level test that is the promotion.
		if transcriptionUnavailable(canvasAwaitsOCR, hasTextAnswer, assistantReady) {
			canvasUnreadable = true
		}

		// A canvas the OCR cannot transcribe must not lower the student's score.
		// It is left ungraded instead of being evaluated as arbitrary text.
		ungraded := canvasUnreadable || statementMediaNeedsReview
		if canvasUnreadable {
			answerText = "UNREADABLE"
			aiFeedback = unreadableCanvasFeedback(teacherGrades)
		} else if statementMediaNeedsReview {
			aiFeedback = statementMediaFeedback(teacherGrades)
		}
		var aiSuggestion *bool

		if statementMediaNeedsReview {
			// The answer is deliberately not auto-graded: see hasStatementMedia.
		} else if ok && ex.Type == exerciseTypeAttachment {
			if !hasAttachment {
				// Nothing was uploaded: that is simply an unanswered exercise.
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
			// Blanks are graded by exact comparison against the expected
			// placement. Sending them to the assistant would spend a call on a
			// settled answer and risk it contradicting the exact match.
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

		// The teacher's queue only ever holds level test answers.
		needsTeacherReview := ungraded && teacherGrades

		score := 0.0
		switch {
		case ungraded:
			// Nobody produced a verdict: drop it from the denominator so an
			// unread answer cannot fail the student on its own.
			total--
			if needsTeacherReview {
				hasPendingReview = true
			}
		case isCorrect:
			correct++
			score = 100.0
		}

		// Track per-topic stats
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
			AIIsCorrect:           aiSuggestion,
		})

		if createErr != nil {
			log.Printf("[practice_submit] could not persist attempt student_id=%s sheet_id=%s err=%v", studentID, sheetID, createErr)
			if ps.SheetType == sheetTypeLevelTest {
				persistenceErr = createErr
				break
			}
		}

		if attempt.CanvasData != "" && attemptID != "" {
			app.Repositories.StudentAttempt.SaveCanvasWork(ctx, attemptID, attempt.CanvasData)
		}

		// Practice sheets teach from their detailed result. A level test must not
		// become an answer oracle, so it deliberately returns no verdict, answer
		// or assistant feedback for an individual exercise.
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

	if ps.SheetType == sheetTypeLevelTest && persistenceErr != nil {
		// Detached from the request: the usual reason a submission fails is that
		// this very context expired, and running the recovery through it meant
		// both queries failed instantly and the student stayed locked out — the
		// exact case the compensation exists for. WithoutCancel keeps the
		// request's values (tracing, logging) and drops only the cancellation.
		cleanupCtx, cancelCleanup := context.WithTimeout(context.WithoutCancel(ctx), 15*time.Second)
		defer cancelCleanup()

		// The claim makes these attempts exclusive to this submission. Delete
		// before release: reversing the order would let a retry race with stale
		// rows from this failed delivery.
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
	// Nothing came back with a verdict: there is no score to act on yet.
	allUngraded := total <= 0 && len(input.Attempts) > 0

	kumon := domain.NewKumonStrategy()

	// For level tests or single-topic sheets, use overall score for level progression
	derivedTopicID := ps.TopicID
	if derivedTopicID == "" && len(topicStats) == 1 {
		for topicID := range topicStats {
			derivedTopicID = topicID
		}
	}

	currentScore := 0.0
	currentLevel := 1
	if ps.SheetType == "level_test" {
		// Level tests advance the course, not the topic. Topic progress can be
		// ahead because of regular practice and must not cause level skipping.
		courseProgress, _ := app.Repositories.CourseProgress.Get(ctx, studentID, ps.CourseID)
		if courseProgress != nil {
			currentLevel = courseProgress.CurrentLevel
		}
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
		// Nothing gradeable came back, so level progression stays where it is
		// rather than acting on a score of zero.
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
		// Some answer still needs a teacher, and this test decides promotion.
		shouldLevelUp = false
		shouldRepeat = false
		nextLevel = currentLevel
		newMastery = currentScore
		recommendation = "Tu prueba quedó esperando la corrección del docente."
	case ps.SheetType == sheetTypeLevelTest:
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
			// Log error but don't fail the submit
			_ = err
		}
	}

	if ps.SheetType == "level_test" && shouldLevelUp {
		app.Repositories.CourseProgress.Upsert(ctx, studentID, ps.CourseID, nextLevel)
	}

	if ps.SheetType == sheetTypeLevelTest {
		resultAIFeedback = ""
	}
	result := toSubmitOutputData(sheetScore, correct, total, newMastery, recommendation, resultAIFeedback, shouldLevelUp, shouldRepeat, nextLevel, exerciseResults)
	result.PendingReview = hasPendingReview
	return &SubmitOutput{Data: result}, nil
}

// validateLevelTestAttempts keeps the server's pass/fail calculation tied to
// the complete test. UI validation alone is bypassable: submitting only one
// easy exercise previously produced a 100% score and promoted the student.
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

// transcriptionUnavailable reports a handwritten answer that never became text.
// Whether the assistant failed or was never reachable, the outcome is the same:
// there is nothing to compare, so the answer has to be left ungraded rather
// than measured against an empty string.
func transcriptionUnavailable(canvasAwaitsOCR, hasTextAnswer, assistantReady bool) bool {
	return canvasAwaitsOCR && !hasTextAnswer && !assistantReady
}

// teacherAssistantConfig resolves the assistant the course owner configured, so
// grading still runs for a student who has no credentials of their own.
func teacherAssistantConfig(ctx context.Context, app *appcontext.Context, teacherID string) assistant.Config {
	cfg := assistant.Config{}
	if strings.TrimSpace(teacherID) == "" {
		return cfg
	}
	teacher, err := app.Repositories.UserProfile.Get(ctx, teacherID)
	if err != nil || teacher == nil {
		return cfg
	}
	cfg.BaseURL = teacher.AssistantBaseURL
	cfg.APIKey = teacher.AssistantAPIKey
	return cfg
}

// teacherGradesSheet says whether a sheet's answers are corrected by a person.
// Only a level test is: it decides promotion, so a human confirms it. A
// practice is formative — it must resolve on submit, so it is never queued for
// a teacher and never left hanging on one.
func teacherGradesSheet(sheetType string) bool {
	return sheetType == sheetTypeLevelTest
}

// unreadableCanvasFeedback explains a handwritten answer the OCR could not
// transcribe. Only a level test can point the student at a correction: on a
// practice nobody is going to look at it, so asking them to request one would
// send them after something the app no longer offers.
func unreadableCanvasFeedback(teacherGrades bool) string {
	if teacherGrades {
		return "No pudimos leer tu respuesta escrita, así que la va a corregir el docente."
	}
	return "No pudimos leer tu respuesta escrita, así que no cuenta en tu puntaje. Intentá escribirla más clara la próxima vez."
}

// statementMediaFeedback explains why an answer to an exercise with image or
// audio in its statement was not auto-graded, without promising a correction
// that only a level test actually gets.
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

// calcStreak counts consecutive days of practice in the student's own zone.
// Measuring in UTC broke it for a UTC-3 audience: see domain.StudentLocation.
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
