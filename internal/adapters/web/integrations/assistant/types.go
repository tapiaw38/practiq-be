package assistant

type Config struct {
	BaseURL string
	APIKey  string
}

type EvaluationResult struct {
	IsCorrect bool
	Feedback  string
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
