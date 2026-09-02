package domain

import "time"

type Notebook struct {
	ID          string
	CourseID    string
	TeacherID   string
	Title       string
	Description string
	Level       int
	Pages       []NotebookPage
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type NotebookPage struct {
	ID          string
	NotebookID  string
	PageNumber  int
	Title       string
	ContentType string // "canvas" | "text"
	ContentData string // teacher content (base64 PNG or text)
	// StatementText is the page transcribed to text when ContentData is an
	// image, so grading compares against words instead of re-reading the
	// picture on every submission.
	StatementText string
	// StatementVerified records that a teacher looked at StatementText. Until
	// then it is an unchecked transcription and a verdict built on it can only
	// be a suggestion.
	StatementVerified bool
	Instructions      string
	Submission        *NotebookSubmission // student's work, if loaded
	CreatedAt         time.Time
}

type NotebookSubmission struct {
	ID                 string
	PageID             string
	StudentID          string
	CanvasData         string
	AnswerText         string
	AIRecognizedText   string
	AIIsCorrect        *bool
	AIFeedback         string
	AIReviewedAt       *time.Time
	NeedsTeacherReview bool
	TeacherIsCorrect   *bool
	TeacherFeedback    string
	TeacherReviewedAt  *time.Time
	SubmittedAt        time.Time
	UpdatedAt          time.Time
}

type NotebookSubmissionFull struct {
	NotebookSubmission
	StudentName   string
	StudentEmail  string
	NotebookID    string
	NotebookTitle string
	PageTitle     string
	PageNumber    int
	CourseID      string
	TeacherID     string
}
