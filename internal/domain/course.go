package domain

import "time"

const (
	CourseStatusDraft = "draft"

	CourseStatusPublished = "published"

	CourseStatusArchived = "archived"
)

func (c Course) VisibleToStudents() bool {
	return c.Status == CourseStatusPublished || c.Status == CourseStatusArchived
}

func (c Course) AcceptsWork() bool {
	return c.Status == CourseStatusPublished
}

func (c Course) AcceptsEnrollment() bool {
	return c.Status == CourseStatusPublished
}

type Course struct {
	ID        string
	TeacherID string

	SchoolID    string
	GradeID     string
	GradeName   string
	GradeTheme  string
	SubjectID   string
	SubjectName string
	Title       string
	Description string
	Level       string
	Subject     string
	Status      string
	CreatedAt   time.Time
}

type CourseCuriosities struct {
	CourseID    string
	Curiosities []string
}

type CourseDashboardSummary struct {
	CourseID string

	SchoolID       string
	SchoolName     string
	Title          string
	Subject        string
	GradeName      string
	PracticeSheets int
	LevelTests     int
	Notebooks      int
	CurrentLevel   int

	TopicIDs []string
}
