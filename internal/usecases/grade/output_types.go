package grade

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	OperationResultData struct {
		Message string `json:"message"`
	}

	GradeData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		CreatedBy   string `json:"created_by"`
		CreatedAt   string `json:"created_at"`
	}

	GradeMemberData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		ProfileType string `json:"profile_type"`
	}
)

func toGradeData(grade domain.Grade) GradeData {
	return GradeData{
		ID:          grade.ID,
		Name:        grade.Name,
		Description: grade.Description,
		CreatedBy:   grade.CreatedBy,
		CreatedAt:   grade.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toOperationResultData(result domain.OperationResult) OperationResultData {
	return OperationResultData{Message: result.Message}
}

func toGradeMemberData(user domain.UserProfile) GradeMemberData {
	return GradeMemberData{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		ProfileType: user.ProfileType,
	}
}
