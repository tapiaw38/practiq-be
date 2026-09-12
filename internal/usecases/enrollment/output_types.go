package enrollment

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	OperationResultData struct {
		Message string `json:"message"`
	}

	StudentData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Email       string `json:"email"`
		ProfileType string `json:"profile_type"`
		CreatedAt   string `json:"created_at"`
	}
)

func toStudentData(p domain.UserProfile, name, email string) StudentData {
	return StudentData{
		ID:          p.ID,
		Name:        name,
		Email:       email,
		ProfileType: p.ProfileType,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toOperationResultData(result domain.OperationResult) OperationResultData {
	return OperationResultData{Message: result.Message}
}
