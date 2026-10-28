package domain

import "time"

type UserProfile struct {
	ID             string
	ProfileType    string
	AcademicStatus string

	Timezone    string
	UITheme     string
	ColorScheme string

	AvatarSeed string
	CreatedAt  time.Time
}
