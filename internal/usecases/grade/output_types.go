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
		VisualTheme string `json:"visual_theme"`
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
		VisualTheme: grade.VisualTheme,
		CreatedBy:   grade.CreatedBy,
		CreatedAt:   grade.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toOperationResultData(result domain.OperationResult) OperationResultData {
	return OperationResultData{Message: result.Message}
}

func toGradeMemberData(user domain.UserProfile, name, email string) GradeMemberData {
	return GradeMemberData{
		ID:          user.ID,
		Name:        name,
		Email:       email,
		ProfileType: user.ProfileType,
	}
}
