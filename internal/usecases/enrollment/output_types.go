package enrollment

import "github.com/tapiaw38/practiq-be/internal/domain"

type StudentData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	ProfileType string `json:"profile_type"`
	CreatedAt   string `json:"created_at"`
}

type EnrollOutput struct {
	Message string `json:"message"`
}

type StudentsListOutput struct {
	Data []StudentData `json:"data"`
}

func toStudentData(p domain.UserProfile) StudentData {
	return StudentData{
		ID:          p.ID,
		Name:        p.Name,
		Email:       p.Email,
		ProfileType: p.ProfileType,
		CreatedAt:   p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
