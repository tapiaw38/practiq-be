package subject

import "github.com/tapiaw38/practiq-be/internal/domain"

type SubjectData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   string `json:"created_by"`
	CreatedAt   string `json:"created_at"`
}

type SubjectOutput struct {
	Data SubjectData `json:"data"`
}

type SubjectListOutput struct {
	Data []SubjectData `json:"data"`
}

func toSubjectData(subject domain.Subject) SubjectData {
	return SubjectData{
		ID:          subject.ID,
		Name:        subject.Name,
		Description: subject.Description,
		CreatedBy:   subject.CreatedBy,
		CreatedAt:   subject.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
