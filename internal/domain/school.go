package domain

import (
	"strings"
	"time"
)

const (
	SchoolKindPersonal    = "personal"
	SchoolKindInstitution = "institution"
)

const (
	SchoolBillingSubscription = "subscription"
	SchoolBillingDirect       = "direct"
)

const (
	SchoolStatusActive    = "active"
	SchoolStatusSuspended = "suspended"
	SchoolStatusClosed    = "closed"
)

const (
	SchoolRoleAdmin   = "admin"
	SchoolRoleTeacher = "teacher"
	SchoolRoleStudent = "student"
)

const PlaceholderSchoolName = "Mi escuela"

type (
	School struct {
		ID          string
		Name        string
		Kind        string
		Billing     string
		Status      string
		CreatedBy   string
		CreatedAt   time.Time
		ClosedAt    *time.Time
		ClosedBy    string
		CloseReason string
	}

	StudentActivity struct {
		UserID        string
		LastPracticed *time.Time
	}

	SchoolMember struct {
		SchoolID string
		UserID   string
		Role     string
		Active   bool
	}
)

func PersonalSchoolName(teacherName string) string {
	name := strings.TrimSpace(teacherName)
	if name == "" {
		return PlaceholderSchoolName
	}
	return "Mi escuela de " + name
}
