package courselevel

import (
	"context"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	practiceSheetRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/practice_sheet"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type (
	GetUsecase interface {
		Execute(ctx context.Context, requesterID string, isAdmin bool, courseID string) (*GetOutput, apperrors.ApplicationError)
	}

	getUsecase struct {
		contextFactory appcontext.Factory
	}

	GetOutput struct {
		CurrentLevel int         `json:"current_level"`
		Levels       []LevelData `json:"levels"`
	}
)

func NewGetUsecase(contextFactory appcontext.Factory) GetUsecase {
	return &getUsecase{contextFactory: contextFactory}
}

func (u *getUsecase) Execute(ctx context.Context, requesterID string, isAdmin bool, courseID string) (*GetOutput, apperrors.ApplicationError) {
	app := u.contextFactory()
	if appErr := requesterCanReadCourse(ctx, app, requesterID, isAdmin, courseID); appErr != nil {
		return nil, appErr
	}

	currentLevel := 1
	if requesterID != "" {
		cp, err := app.Repositories.CourseProgress.Get(ctx, requesterID, courseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
		if cp != nil {
			currentLevel = cp.CurrentLevel
		}
	}

	sheets, err := app.Repositories.PracticeSheet.List(ctx, practiceSheetRepo.ListFilter{
		CourseID: courseID,
	})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	notebooks, err := app.Repositories.Notebook.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	maxLevel := currentLevel
	for _, s := range sheets {
		if s.Level > maxLevel {
			maxLevel = s.Level
		}
	}
	for _, nb := range notebooks {
		if nb.Level > maxLevel {
			maxLevel = nb.Level
		}
	}
	if maxLevel <= currentLevel {
		maxLevel = currentLevel + 1
	}

	type levelBucket struct {
		practices []SheetData
		levelTest *SheetData
		notebooks []NotebookData
	}
	buckets := make(map[int]*levelBucket)
	for i := 1; i <= maxLevel; i++ {
		buckets[i] = &levelBucket{
			practices: []SheetData{},
			notebooks: []NotebookData{},
		}
	}

	for _, s := range sheets {
		b := buckets[s.Level]
		if b == nil {
			continue
		}
		sd := toSheetData(s)
		if s.SheetType == "level_test" {
			b.levelTest = &sd
		} else {
			b.practices = append(b.practices, sd)
		}
	}

	for _, nb := range notebooks {
		b := buckets[nb.Level]
		if b == nil {
			continue
		}
		b.notebooks = append(b.notebooks, toNotebookData(nb))
	}

	levels := make([]LevelData, 0, maxLevel)
	for i := 1; i <= maxLevel; i++ {
		b := buckets[i]
		levels = append(levels, LevelData{
			Level:     i,
			Unlocked:  i <= currentLevel,
			Practices: b.practices,
			LevelTest: b.levelTest,
			Notebooks: b.notebooks,
		})
	}

	return &GetOutput{
		CurrentLevel: currentLevel,
		Levels:       levels,
	}, nil
}

func requesterCanReadCourse(ctx context.Context, app *appcontext.Context, requesterID string, isAdmin bool, courseID string) apperrors.ApplicationError {
	if isAdmin {
		return nil
	}
	course, err := app.Repositories.Course.Get(ctx, courseID)
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	if course == nil {
		return apperrors.NewNotFoundError("course not found")
	}
	if course.TeacherID == requesterID {
		return nil
	}
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{StudentID: requesterID})
	if err != nil {
		return apperrors.NewApplicationError(mappings.InternalServerError, err)
	}
	for _, course := range courses {
		if course.ID == courseID {
			return nil
		}
	}
	return apperrors.NewForbiddenError()
}
