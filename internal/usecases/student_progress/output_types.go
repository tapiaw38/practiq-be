package studentprogress

import (
	"time"

	"github.com/tapiaw38/practiq-be/internal/domain"
)

type ProgressData struct {
	TopicID         string  `json:"topic_id"`
	TopicTitle      string  `json:"topic_title"`
	StrategyID      string  `json:"strategy_id"`
	MasteryScore    float64 `json:"mastery_score"`
	CurrentLevel    int     `json:"current_level"`
	TotalAttempts   int     `json:"total_attempts"`
	CorrectAttempts int     `json:"correct_attempts"`
	StreakDays      int     `json:"streak_days"`
	LastPracticedAt string  `json:"last_practiced_at"`
}

type ProgressListOutput struct {
	Data []ProgressData `json:"data"`
}

// effectiveStreak returns 0 when the streak is broken: the stored value is
// only valid if the student practiced today or yesterday.
func effectiveStreak(p domain.StudentTopicProgress) int {
	if p.LastPracticedAt == nil {
		return 0
	}
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	yesterday := today.AddDate(0, 0, -1)
	last := p.LastPracticedAt.In(now.Location())
	lastDay := time.Date(last.Year(), last.Month(), last.Day(), 0, 0, 0, 0, now.Location())
	if lastDay.Equal(today) || lastDay.Equal(yesterday) {
		return p.StreakDays
	}
	return 0
}

func toProgressData(p domain.StudentTopicProgress) ProgressData {
	lastPracticed := ""
	if p.LastPracticedAt != nil {
		lastPracticed = p.LastPracticedAt.Format("2006-01-02T15:04:05Z")
	}
	return ProgressData{
		TopicID:         p.TopicID,
		TopicTitle:      p.TopicTitle,
		StrategyID:      p.StrategyID,
		MasteryScore:    p.MasteryScore,
		CurrentLevel:    p.CurrentLevel,
		TotalAttempts:   p.TotalAttempts,
		CorrectAttempts: p.CorrectAttempts,
		StreakDays:      effectiveStreak(p),
		LastPracticedAt: lastPracticed,
	}
}
