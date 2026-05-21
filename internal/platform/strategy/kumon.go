package strategy

import "math"

type MasteryInput struct {
	TotalAttempts    int
	CorrectAttempts  int
	HintsUsed        int
	TimeSpentSeconds int
	CurrentScore     float64
}

type Recommendation struct {
	Action    string
	Message   string
	NextLevel int
}

type LearningStrategyService interface {
	CalculateMasteryScore(input MasteryInput) float64
	ShouldLevelUp(masteryScore float64) bool
	ShouldRepeatTopic(masteryScore float64) bool
	GenerateNextPracticeRecommendation(masteryScore float64, currentLevel int) *Recommendation
}

type kumonStrategy struct{}

func NewKumonStrategy() LearningStrategyService {
	return &kumonStrategy{}
}

func (k *kumonStrategy) CalculateMasteryScore(input MasteryInput) float64 {
	if input.TotalAttempts == 0 {
		return input.CurrentScore
	}

	baseScore := float64(input.CorrectAttempts) / float64(input.TotalAttempts) * 100

	// Deduct 2 points per hint used
	hintPenalty := float64(input.HintsUsed) * 2
	baseScore = math.Max(0, baseScore-hintPenalty)

	// Bonus for perfect score with no hints
	if input.CorrectAttempts == input.TotalAttempts && input.HintsUsed == 0 {
		baseScore = math.Min(100, baseScore+5)
	}

	// Weighted average: 70% new score + 30% existing
	newScore := baseScore*0.7 + input.CurrentScore*0.3
	return math.Round(newScore*100) / 100
}

func (k *kumonStrategy) ShouldLevelUp(masteryScore float64) bool {
	return masteryScore >= 90
}

func (k *kumonStrategy) ShouldRepeatTopic(masteryScore float64) bool {
	return masteryScore < 70
}

func (k *kumonStrategy) GenerateNextPracticeRecommendation(masteryScore float64, currentLevel int) *Recommendation {
	if k.ShouldLevelUp(masteryScore) {
		return &Recommendation{
			Action:    "level_up",
			Message:   "Excelente dominio! Puedes avanzar al siguiente nivel.",
			NextLevel: currentLevel + 1,
		}
	}

	if k.ShouldRepeatTopic(masteryScore) {
		return &Recommendation{
			Action:    "repeat",
			Message:   "Necesitas practicar más este tema antes de continuar.",
			NextLevel: currentLevel,
		}
	}

	return &Recommendation{
		Action:    "continue",
		Message:   "Vas por buen camino. Sigue practicando para dominar el tema.",
		NextLevel: currentLevel,
	}
}
