package topic

import "github.com/tapiaw38/practiq-be/internal/domain"

type TopicData struct {
	ID          string `json:"id"`
	CourseID    string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderIndex  int    `json:"order_index"`
	CreatedAt   string `json:"created_at"`
}

type TopicOutput struct {
	Data TopicData `json:"data"`
}

type TopicListOutput struct {
	Data []TopicData `json:"data"`
}

func toTopicData(t domain.Topic) TopicData {
	return TopicData{
		ID:          t.ID,
		CourseID:    t.CourseID,
		Title:       t.Title,
		Description: t.Description,
		OrderIndex:  t.OrderIndex,
		CreatedAt:   t.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
