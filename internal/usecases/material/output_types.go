package material

import "github.com/tapiaw38/practiq-be/internal/domain"

type MaterialData struct {
	ID            string `json:"id"`
	CourseID      string `json:"course_id"`
	TeacherID     string `json:"teacher_id"`
	Title         string `json:"title"`
	Type          string `json:"type"`
	FileURL       string `json:"file_url,omitempty"`
	ExtractedText string `json:"extracted_text,omitempty"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type MaterialOutput struct {
	Data MaterialData `json:"data"`
}

type MaterialListOutput struct {
	Data []MaterialData `json:"data"`
}

func toMaterialData(m domain.Material) MaterialData {
	return MaterialData{
		ID:            m.ID,
		CourseID:      m.CourseID,
		TeacherID:     m.TeacherID,
		Title:         m.Title,
		Type:          m.Type,
		FileURL:       m.FileURL,
		ExtractedText: m.ExtractedText,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
