package course

import "github.com/tapiaw38/practiq-be/internal/domain"

type CourseData struct {
	ID          string `json:"id"`
	TeacherID   string `json:"teacher_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	Subject     string `json:"subject"`
	CreatedAt   string `json:"created_at"`
}

type CourseOutput struct {
	Data CourseData `json:"data"`
}

type CourseListOutput struct {
	Data []CourseData `json:"data"`
}

func toCourseData(c domain.Course) CourseData {
	return CourseData{
		ID:          c.ID,
		TeacherID:   c.TeacherID,
		Title:       c.Title,
		Description: c.Description,
		Level:       c.Level,
		Subject:     c.Subject,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
