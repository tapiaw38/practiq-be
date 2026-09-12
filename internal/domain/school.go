package domain

import (
	"strings"
	"time"
)

// School kinds. What an admin may do depends on this rather than on their role:
// a teacher administering their own school and an administrator running an
// institution are both `admin`, and only one of them may add teachers.
const (
	SchoolKindPersonal    = "personal"
	SchoolKindInstitution = "institution"
)

// Billing modes. Direct means invoiced outside the product: no entitlement is
// consulted and no student limit applies.
const (
	SchoolBillingSubscription = "subscription"
	SchoolBillingDirect       = "direct"
)

// School lifecycle. Closed schools retain every record but disappear from
// members' scopes. Suspended is reserved for temporary operational pauses.
const (
	SchoolStatusActive    = "active"
	SchoolStatusSuspended = "suspended"
	SchoolStatusClosed    = "closed"
)

// Roles inside a school. They belong to the membership, not to the person, so
// the same teacher can administer their own school and teach at an institution.
const (
	SchoolRoleAdmin   = "admin"
	SchoolRoleTeacher = "teacher"
	SchoolRoleStudent = "student"
)

// PlaceholderSchoolName is what the schools migration had to use. Teacher names
// live in auth-api-be, so the SQL could not know them; the application replaces
// this the first time it can resolve the real one.
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

	SchoolMember struct {
		SchoolID string
		UserID   string
		Role     string
		Active   bool
	}
)

// PersonalSchoolName is what a teacher's own school is called before they
// rename it. Falling back to the placeholder rather than to a half name: a
// name is shown to the teacher, and an empty suffix reads as a bug.
func PersonalSchoolName(teacherName string) string {
	name := strings.TrimSpace(teacherName)
	if name == "" {
		return PlaceholderSchoolName
	}
	return "Mi escuela de " + name
}
