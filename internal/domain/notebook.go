package domain

import "time"

type Notebook struct {
	ID          string
	CourseID    string
	TeacherID   string
	Title       string
	Description string
	Pages       []NotebookPage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NotebookPage struct {
	ID           string
	NotebookID   string
	PageNumber   int
	Title        string
	ContentType  string // "canvas" | "text"
	ContentData  string // teacher content (base64 PNG or text)
	Instructions string
	Submission   *NotebookSubmission // student's work, if loaded
	CreatedAt    time.Time
}

type NotebookSubmission struct {
	ID          string
	PageID      string
	StudentID   string
	CanvasData  string
	AnswerText  string
	SubmittedAt time.Time
	UpdatedAt   time.Time
}
