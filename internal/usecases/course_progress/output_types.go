package courseprogress

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	CourseProgressData struct {
		ID           string `json:"id,omitempty"`
		StudentID    string `json:"student_id"`
		CourseID     string `json:"course_id"`
		CurrentLevel int    `json:"current_level"`
		UpdatedAt    string `json:"updated_at,omitempty"`
	}
)

func toProgressData(p domain.StudentCourseProgress) CourseProgressData {
	updatedAt := ""
	if !p.UpdatedAt.IsZero() {
		updatedAt = p.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}
	return CourseProgressData{
		ID:           p.ID,
		StudentID:    p.StudentID,
		CourseID:     p.CourseID,
		CurrentLevel: p.CurrentLevel,
		UpdatedAt:    updatedAt,
	}
}
