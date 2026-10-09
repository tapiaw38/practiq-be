package studentprogress

import (
	"context"
	"strings"

	courseRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/course"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/integrations/authapi"
	"github.com/tapiaw38/practiq-be/internal/platform/appcontext"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	"github.com/tapiaw38/practiq-be/internal/platform/errors/mappings"
	"github.com/tapiaw38/practiq-be/internal/platform/identity"
)

const leaderboardTop = 10

type (
	GetCourseLeaderboardUsecase interface {
		Execute(ctx context.Context, studentID, courseID, bearerToken string) (*GetCourseLeaderboardOutput, apperrors.ApplicationError)
	}

	getCourseLeaderboardUsecase struct {
		contextFactory appcontext.Factory
	}

	LeaderboardEntryData struct {
		Name string `json:"name"`

		AvatarSeed string `json:"avatar_seed,omitempty"`
		TotalXP    int    `json:"total_xp"`
		Position   int    `json:"position"`
		IsMe       bool   `json:"is_me,omitempty"`
	}

	GetCourseLeaderboardOutput struct {
		Data []LeaderboardEntryData `json:"data"`

		Me *LeaderboardEntryData `json:"me,omitempty"`
	}
)

func NewGetCourseLeaderboardUsecase(contextFactory appcontext.Factory) GetCourseLeaderboardUsecase {
	return &getCourseLeaderboardUsecase{contextFactory: contextFactory}
}

func (u *getCourseLeaderboardUsecase) Execute(ctx context.Context, studentID, courseID, bearerToken string) (*GetCourseLeaderboardOutput, apperrors.ApplicationError) {
	app := u.contextFactory()

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

	ids := make([]string, 0, len(entries))
	for _, entry := range entries {
		ids = append(ids, entry.StudentID)
	}
	names, appErr := identity.Names(ctx, app.Integrations.AuthAPI, bearerToken, ids)
	if appErr != nil {
		return nil, appErr
	}

	output := GetCourseLeaderboardOutput{Data: make([]LeaderboardEntryData, 0, len(entries))}
	for i, entry := range entries {
		item := LeaderboardEntryData{
			Name:       shortDisplayName(names[entry.StudentID]),
			AvatarSeed: entry.AvatarSeed,
			TotalXP:    entry.TotalXP,
			Position:   entry.Position,
			IsMe:       entry.StudentID == studentID,
		}

		if item.IsMe && i >= leaderboardTop {
			me := item
			output.Me = &me
			continue
		}
		output.Data = append(output.Data, item)
	}

	return &output, nil
}

func shortDisplayName(info authapi.UserInfo) string {
	first := strings.TrimSpace(info.FirstName)
	last := strings.TrimSpace(info.LastName)
	if first == "" && last == "" {
		return "Alumno"
	}
	if first == "" {
		first, last = last, ""
	}
	if last == "" {
		return first
	}
	return first + " " + strings.ToUpper(string([]rune(last)[:1])) + "."
}
