package courselevel

import (
	"context"

	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

type GetUsecase interface {
	Execute(ctx context.Context, courseID, studentID string) (*CourseLevelsOutput, apperrors.ApplicationError)
}

type getUsecase struct {
	factory appcontext.Factory
}

func NewGetUsecase(factory appcontext.Factory) GetUsecase {
	return &getUsecase{factory: factory}
}

func (u *getUsecase) Execute(ctx context.Context, courseID, studentID string) (*CourseLevelsOutput, apperrors.ApplicationError) {
	app := u.factory()

	// Student's current level for this course
	currentLevel := 1
	if studentID != "" {
		cp, err := app.Repositories.CourseProgress.Get(ctx, studentID, courseID)
		if err != nil {
			return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
		}
		if cp != nil {
			currentLevel = cp.CurrentLevel
		}
	}

	// All sheets for the course
	sheets, err := app.Repositories.PracticeSheet.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	// All notebooks for the course
	notebooks, err := app.Repositories.Notebook.List(ctx, courseID)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.InternalServerError, err)
	}

	// Determine max level across all content
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
	// Always show at least one level ahead (locked preview)
	if maxLevel <= currentLevel {
		maxLevel = currentLevel + 1
	}

	// Build level map
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

	return &CourseLevelsOutput{
		CurrentLevel: currentLevel,
		Levels:       levels,
	}, nil
}
