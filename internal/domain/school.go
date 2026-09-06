package domain

import "time"

type School struct {
	ID        string
	Name      string
	Slug      string
	IsActive  bool
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SchoolMembership struct {
	SchoolID       string
	UserID         string
	MembershipRole string
	ProfileType    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
