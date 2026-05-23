package course

import "github.com/tapiaw38/practiq-be/internal/domain"

type CourseData struct {
	ID          string `json:"id"`
	TeacherID   string `json:"teacher_id"`
	GradeID     string `json:"grade_id"`
	GradeName   string `json:"grade_name"`
	SubjectID   string `json:"subject_id"`
	SubjectName string `json:"subject_name"`
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
		GradeID:     c.GradeID,
		GradeName:   c.GradeName,
		SubjectID:   c.SubjectID,
		SubjectName: c.SubjectName,
		Title:       c.Title,
		Description: c.Description,
		Level:       c.Level,
		Subject:     c.Subject,
		CreatedAt:   c.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
