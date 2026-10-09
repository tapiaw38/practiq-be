package assistant

import "errors"

type Config struct {
	BaseURL string
	APIKey  string
}

type EvaluationResult struct {
	IsCorrect bool
	Feedback  string
}

const (
	AttachmentKindAudio    = "audio"
	AttachmentKindImage    = "image"
	AttachmentKindPDF      = "pdf"
	AttachmentKindDocument = "doc"
)

var ErrAttachmentNotEvaluable = errors.New("attachment kind cannot be evaluated by the assistant")

const UnreadableFeedback = "UNREADABLE"

type AttachmentEvaluationInput struct {
	Question      string
	CorrectAnswer string
	GradeName     string
	Kind          string
	Filename      string
	ContentType   string
	Content       []byte
}

type ProxyResponse struct {
	StatusCode  int
	ContentType string
	Body        []byte
}

type createConversationRequest struct {
	Title     string `json:"title"`
	IsSandbox bool   `json:"is_sandbox"`
}

type createConversationResponse struct {
	Data struct {
		ID string `json:"id"`
	} `json:"data"`
}

type messageResponse struct {
	Data []struct {
		Content string `json:"content"`
		Sender  string `json:"sender"`
	} `json:"data"`
}
