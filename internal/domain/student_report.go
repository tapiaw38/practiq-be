package domain

import "time"

type StudentReportFilter struct {
	StudentID string
	CourseID  string
	From      *time.Time
	To        *time.Time
}

type StudentReportData struct {
	Student         UserProfile
	GeneratedAt     time.Time
	Period          ReportPeriod
	Summary         ReportSummary
	TopicProgress   []StudentTopicProgress
	CourseProgress  []CourseProgressItem
	DailyAttempts   []DailyAttemptCount
	RecentAttempts  []StudentAttempt
}

type ReportPeriod struct {
	From *time.Time
	To   *time.Time
}

type ReportSummary struct {
	TopicsPracticed   int
	AverageMastery    float64
	TotalAttempts     int
	CorrectAttempts   int
	AccuracyRate      float64
	CurrentStreak     int
}

type CourseProgressItem struct {
	CourseID      string
	CourseTitle   string
	CurrentLevel  int
	TopicCount    int
	AverageMastery float64
	LastActivity  *time.Time
}

type DailyAttemptCount struct {
	Date    time.Time
	Total   int
	Correct int
}
