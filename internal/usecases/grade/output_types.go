package grade

import "github.com/tapiaw38/practiq-be/internal/domain"

type GradeData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}

type GradeOutput struct {
	Data GradeData `json:"data"`
}

type GradeListOutput struct {
	Data []GradeData `json:"data"`
}

type GradeMemberData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ProfileType string `json:"profile_type"`
}

type GradeMembersOutput struct {
	Data []GradeMemberData `json:"data"`
}

type AssignMemberOutput struct {
	Message string `json:"message"`
}

type RemoveMemberOutput struct {
	Message string `json:"message"`
}

func toGradeData(grade domain.Grade) GradeData {
	return GradeData{
		ID:          grade.ID,
		Name:        grade.Name,
		Description: grade.Description,
		CreatedBy:   grade.CreatedBy,
		CreatedAt:   grade.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toGradeMemberData(user domain.UserProfile) GradeMemberData {
	return GradeMemberData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		ProfileType: user.ProfileType,
	}
}
