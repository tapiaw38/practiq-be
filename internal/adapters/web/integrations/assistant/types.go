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

// Attachment kinds supported by Gillie's conversation multipart channels.
const (
	AttachmentKindAudio    = "audio"
	AttachmentKindImage    = "image"
	AttachmentKindPDF      = "pdf"
	AttachmentKindDocument = "doc"
)

// ErrAttachmentNotEvaluable means the assistant cannot grade this kind of file
// and a human has to.
var ErrAttachmentNotEvaluable = errors.New("attachment kind cannot be evaluated by the assistant")

// UnreadableFeedback is what the assistant answers when it received the file
// but could not make sense of it.
const UnreadableFeedback = "UNREADABLE"

type AttachmentEvaluationInput struct {
	Question      string
	CorrectAnswer string
	GradeName     string
	Kind          string
	Filename      string
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
