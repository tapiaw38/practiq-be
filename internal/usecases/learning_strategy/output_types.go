package learningstrategy

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	LearningStrategyData struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      string `json:"status"`
		CreatedAt   string `json:"created_at"`
	}

	CourseLearningStrategyData struct {
		ID                  string `json:"id"`
		CourseID            string `json:"course_id"`
		StrategyID          string `json:"strategy_id"`
		IsDefault           bool   `json:"is_default"`
		Config              string `json:"config"`
		CreatedAt           string `json:"created_at"`
		StrategyName        string `json:"strategy_name"`
		StrategyCode        string `json:"strategy_code"`
		StrategyDescription string `json:"strategy_description"`
	}
)

func toStrategyData(s domain.LearningStrategy) LearningStrategyData {
	return LearningStrategyData{
		ID:          s.ID,
		Name:        s.Name,
		Code:        s.Code,
		Description: s.Description,
		Status:      s.Status,
		CreatedAt:   s.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toCourseStrategyData(cls domain.CourseLearningStrategy) CourseLearningStrategyData {
	return CourseLearningStrategyData{
		ID:                  cls.ID,
		CourseID:            cls.CourseID,
		StrategyID:          cls.StrategyID,
		IsDefault:           cls.IsDefault,
		Config:              cls.Config,
		CreatedAt:           cls.CreatedAt.Format("2006-01-02T15:04:05Z"),
		StrategyName:        cls.StrategyName,
		StrategyCode:        cls.StrategyCode,
		StrategyDescription: cls.StrategyDescription,
	}
}
