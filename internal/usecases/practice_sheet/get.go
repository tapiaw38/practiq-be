package practicesheet

import (
	"context"
	"log"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(context.Context, string, bool, string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		Data PracticeSheetData `json:"data"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, requesterID string, isSuperAdmin bool, id string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	ps, err := app.Repositories.PracticeSheet.Get(ctx, id)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
	}
	if ps == nil {
		return nil, apperrors.NewApplicationError(mappings.PracticeSheetNotFoundError, nil)
	}
	if appErr := requesterCanReadCourse(ctx, app, requesterID, isSuperAdmin, ps.CourseID); appErr != nil {
		return nil, appErr
	}
	if appErr := ensureSheetIsOpen(ctx, app, ps, requesterID, isSuperAdmin); appErr != nil {
		return nil, appErr
	}
	includeTeacherData, appErr := requesterCanViewTeacherData(ctx, app, requesterID, isSuperAdmin, ps.CourseID)
	if appErr != nil {
		return nil, appErr
	}
	// Only student practice sheets feed "Continuar práctica". Level tests have
	// a separate route and timing rules, so they must never replace this state.
	if ps.SheetType != sheetTypeLevelTest && !includeTeacherData && requesterID != "" {
		if err := app.Repositories.StudentPracticeState.MarkOpened(ctx, requesterID, ps.ID); err != nil {
			log.Printf("[practice_sheet] could not save last opened sheet_id=%s student_id=%s err=%v", ps.ID, requesterID, err)
		}
	}

	data := toSheetData(app, *ps, includeTeacherData)

	// Opening the test is what starts a student's clock, so there is no
	// separate "begin" step to forget or to skip by going straight to submit.
	// A teacher reading their own sheet is not taking it, so they start none.
	if ps.SheetType == sheetTypeLevelTest && !includeTeacherData && requesterID != "" {
		if ps.TimeLimitMinutes != nil {
			if err := app.Repositories.StudentAttempt.MarkLevelTestStarted(ctx, requesterID, ps.ID); err != nil {
				// The test still opens: a clock that failed to start is worth
				// less than a student locked out of a test they can see.
				log.Printf("[practice_sheet] could not start the clock sheet_id=%s student_id=%s err=%v", ps.ID, requesterID, err)
			}
		}
		attempts, startedAt, err := app.Repositories.StudentAttempt.LevelTestProgress(ctx, requesterID, ps.ID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.PracticeSheetGetError, err)
		}
		data.AttemptsUsed = attempts
		data.AttemptsAllowed = ps.AttemptsAllowed()
		if deadline := ps.Deadline(startedAt); deadline != nil {
			data.Deadline = deadline.UTC().Format(timeFormat)
		}
	}

	return &GetOutput{Data: data}, nil
}
