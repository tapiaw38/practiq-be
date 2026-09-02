package domain

import "time"

type Course struct {
	ID          string
	TeacherID   string
	GradeID     string
	GradeName   string
	GradeTheme  string
	SubjectID   string
	SubjectName string
	Title       string
	Description string
	Level       string
	Subject     string
	CreatedAt   time.Time
}

type CourseCuriosities struct {
	CourseID    string
	Curiosities []string
}

// CourseDashboardSummary is one row of the student home: what the screen needs
// about a course without fetching the course's sheets, notebooks and levels.
type CourseDashboardSummary struct {
	CourseID       string
	Title          string
	Subject        string
	PracticeSheets int
	LevelTests     int
	Notebooks      int
	CurrentLevel   int
	// TopicIDs lets the home flag topics that need review without fetching
	// the course's practice sheets to find out which topics it covers.
	TopicIDs []string
}
