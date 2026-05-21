package studentprogress

import "github.com/tapiaw38/practiq-be/internal/domain"

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
		StreakDays:      p.StreakDays,
		LastPracticedAt: lastPracticed,
	}
}
