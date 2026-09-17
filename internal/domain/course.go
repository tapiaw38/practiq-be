package domain

import "time"

const (
	// CourseStatusDraft is a course being prepared. Students never see it,
	// not even one already enrolled: nothing in it is ready to be read.
	CourseStatusDraft = "draft"
	// CourseStatusPublished is a course running normally.
	CourseStatusPublished = "published"
	// CourseStatusArchived is a course that has finished. Students already in
	// it keep reading their work and their marks — taking that away the day a
	// teacher tidies up would erase their own history — but nobody new joins
	// and nothing more can be handed in.
	CourseStatusArchived = "archived"
)

// VisibleToStudents reports whether a student may open the course at all.
func (c Course) VisibleToStudents() bool {
	return c.Status == CourseStatusPublished || c.Status == CourseStatusArchived
}

// AcceptsWork reports whether students may still submit. Archived courses are
// readable but closed to new work.
func (c Course) AcceptsWork() bool {
	return c.Status == CourseStatusPublished
}

// AcceptsEnrollment reports whether somebody new may join.
func (c Course) AcceptsEnrollment() bool {
	return c.Status == CourseStatusPublished
}

type Course struct {
	ID        string
	TeacherID string
	// SchoolID is derived from its grade/subject. It is carried at read time so
	// every access path can enforce the school's lifecycle.
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

// CourseDashboardSummary is one row of the student home: what the screen needs
// about a course without fetching the course's sheets, notebooks and levels.
type CourseDashboardSummary struct {
	CourseID string
	// SchoolID and SchoolName let a student who studies at two schools filter
	// their home. Derived from the course's grade or its teacher: courses carry
	// no school of their own, and a second source would be the same fact twice.
	SchoolID       string
	SchoolName     string
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
