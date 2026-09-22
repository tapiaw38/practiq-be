package studentprogress

import (
	"context"
	"strings"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
)

// leaderboardTop is how many places the table shows before it cuts to the
// student's own row.
const leaderboardTop = 10

type (
	GetCourseLeaderboardUsecase interface {
		Execute(ctx context.Context, studentID, courseID string) (*GetCourseLeaderboardOutput, apperrors.ApplicationError)
	}

	getCourseLeaderboardUsecase struct {
		contextFactory appcontext.Factory
	}

	LeaderboardEntryData struct {
		// Shortened for display: classmates recognise each other by first name,
		// and a ranking of minors is no place for a full surname.
		Name     string `json:"name"`
		TotalXP  int    `json:"total_xp"`
		Position int    `json:"position"`
		IsMe     bool   `json:"is_me,omitempty"`
	}

	GetCourseLeaderboardOutput struct {
		Data []LeaderboardEntryData `json:"data"`
		// Absent while the student has not earned anything: they are nowhere in
		// the ranking yet, which is not the same as being last.
		Me *LeaderboardEntryData `json:"me,omitempty"`
	}
)

func NewGetCourseLeaderboardUsecase(contextFactory appcontext.Factory) GetCourseLeaderboardUsecase {
	return &getCourseLeaderboardUsecase{contextFactory: contextFactory}
}

func (u *getCourseLeaderboardUsecase) Execute(ctx context.Context, studentID, courseID string) (*GetCourseLeaderboardOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

	// The ranking exposes classmates, so course membership is the gate: without
	// this, any authenticated user could read any course's table by id.
	courses, err := app.Repositories.Course.List(ctx, courseRepo.ListFilterOptions{StudentID: studentID})
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}
	belongs := false
	for _, course := range courses {
		if course.ID == courseID {
			belongs = true
			break
		}
	}
	if !belongs {
		return nil, apperrors.NewForbiddenError()
	}

	entries, err := app.Repositories.StudentCourseXP.LeaderboardByCourse(ctx, courseID, studentID, leaderboardTop)
	if err != nil {
		return nil, apperrors.NewApplicationError(mappings.ProgressGetError, err)
	}

	output := GetCourseLeaderboardOutput{Data: make([]LeaderboardEntryData, 0, len(entries))}
	for _, entry := range entries {
		item := LeaderboardEntryData{
			Name:     shortDisplayName(entry.Name),
			TotalXP:  entry.TotalXP,
			Position: entry.Position,
			IsMe:     entry.StudentID == studentID,
		}
		if item.IsMe {
			me := item
			output.Me = &me
		}
		// The student's own row is repeated outside the top on purpose; keeping
		// it in Data too would print it twice when they are inside it.
		if item.Position <= leaderboardTop {
			output.Data = append(output.Data, item)
		}
	}

	return &output, nil
}

// shortDisplayName turns "Walter Tapia Gomez" into "Walter T.". A single word
// is left alone, and a profile with no name still gets something to render.
func shortDisplayName(name string) string {
	parts := strings.Fields(name)
	switch len(parts) {
	case 0:
		return "Alumno"
	case 1:
		return parts[0]
	default:
		return parts[0] + " " + strings.ToUpper(string([]rune(parts[1])[:1])) + "."
	}
}
