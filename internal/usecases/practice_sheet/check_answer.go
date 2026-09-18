package practicesheet

import (
	"context"
	"log"
	"strings"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/usecases/assistantcfg"
	"github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type (
	CheckAnswerUsecase interface {
		Execute(ctx context.Context, sheetID, exerciseID, studentID string, input CheckAnswerInput) (*CheckAnswerOutput, apperrors.ApplicationError)
	}

	checkAnswerUsecase struct {
		contextFactory appcontext.Factory
	}

	CheckAnswerInput struct {
		AnswerText string `json:"answer_text"`
	}

	CheckAnswerResult struct {
		IsCorrect bool   `json:"is_correct"`
		Graded    bool   `json:"graded"`
		Feedback  string `json:"feedback,omitempty"`
	}

	CheckAnswerOutput struct {
		Data CheckAnswerResult `json:"data"`
	}
)

func NewCheckAnswerUsecase(contextFactory appcontext.Factory) CheckAnswerUsecase {
	return &checkAnswerUsecase{contextFactory: contextFactory}
}

// checkableExercise reports whether an answer to this exercise can be judged
// from its text alone. Handwriting needs OCR, an attachment needs to be fetched
// and read, and a statement carrying an image or audio can mean something the
// question text does not say — all of which the submit flow does, slowly and
// once. Those keep their single verdict at submit time.
func checkableExercise(ex domain.Exercise) bool {
	return ex.Type != exerciseTypeAttachment && ex.MediaURL() == ""
}

// Execute judges one answer without recording anything. It exists so a practice
// can react while the student is still on the exercise; the submission remains
// the only thing that is stored and scored.
func (u *checkAnswerUsecase) Execute(ctx context.Context, sheetID, exerciseID, studentID string, input CheckAnswerInput) (*CheckAnswerOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, sheetID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}

	// A level test is an assessment. Telling the student mid-test whether an
	// answer is right would turn the single graded attempt into unlimited ones.
	if ps.SheetType == sheetTypeLevelTest {
		return nil, apperrors.NewBadRequestError("a level test cannot be checked answer by answer")
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
	if appErr := ensureSheetIsOpen(ctx, app, ps, studentID, false); appErr != nil {
		return nil, appErr
	}

	var exercise *domain.Exercise
	for _, pse := range ps.Exercises {
		if pse.Exercise.ID == exerciseID {
			found := pse.Exercise
			exercise = &found
			break
		}
	}
	if exercise == nil {
		return nil, apperrors.NewNotFoundError("exercise not found on this sheet")
	}
	if !checkableExercise(*exercise) {
		return nil, apperrors.NewBadRequestError("this exercise is only graded when the practice is submitted")
	}

	answerText := strings.TrimSpace(input.AnswerText)
	if answerText == "" {
		return nil, apperrors.NewBadRequestError("answer_text is required")
	}
	if isDataURIAnswer(answerText) {
		return nil, apperrors.NewBadRequestError("this exercise is only graded when the practice is submitted")
	}

	if exercise.Type == exerciseTypeFillBlanks {
		return &CheckAnswerOutput{Data: CheckAnswerResult{
			IsCorrect: blanksAnswersMatch(answerText, exercise.CorrectAnswer),
			Graded:    true,
		}}, nil
	}

	isCorrect := strings.EqualFold(
		normalizeCanvasAnswer(answerText),
		normalizeCanvasAnswer(exercise.CorrectAnswer),
	)

	assistantCfg := assistantcfg.Resolve(ctx, app)
	if app.Integrations.AssistantGateway == nil || !app.Integrations.AssistantGateway.IsConfigured(assistantCfg) {
		// Without the assistant only an exact match can be trusted. A different
		// wording of the same answer must not be called wrong on the spot, so it
		// is left for the submission to judge.
		return &CheckAnswerOutput{Data: CheckAnswerResult{IsCorrect: isCorrect, Graded: isCorrect}}, nil
	}

	gradeName := ""
	if course, _ := app.Repositories.Course.Get(ctx, ps.CourseID); course != nil {
		gradeName = course.GradeName
	}

	evaluation, aiErr := app.Integrations.AssistantGateway.EvaluatePracticeAnswer(
		ctx, assistantCfg, exercise.Question, exercise.CorrectAnswer, answerText, gradeName,
	)
	if aiErr != nil {
		log.Printf("[practice_check] evaluation unavailable student_id=%s exercise_id=%s err=%v", studentID, exerciseID, aiErr)
		return &CheckAnswerOutput{Data: CheckAnswerResult{IsCorrect: isCorrect, Graded: isCorrect}}, nil
	}

	return &CheckAnswerOutput{Data: CheckAnswerResult{
		IsCorrect: evaluation.IsCorrect,
		Graded:    true,
		Feedback:  evaluation.Feedback,
	}}, nil
}
