package domain

import "time"

type Notebook struct {
	ID          string
	CourseID    string
	TopicID     string
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
	ContentType string
	ContentData string

	StatementText string

	StatementVerified bool
	Instructions      string
	Submission        *NotebookSubmission
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

	Version           int64
	TeacherIsCorrect  *bool
	TeacherFeedback   string
	TeacherReviewedAt *time.Time
	SubmittedAt       time.Time
	UpdatedAt         time.Time
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
