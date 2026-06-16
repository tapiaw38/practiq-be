package assistant

import (
	"context"
	"net/http"
	"strings"
	"time"
)

type Gateway interface {
	IsConfigured(cfg Config) bool
	AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error)
	EvaluatePracticeAnswer(ctx context.Context, cfg Config, question, correctAnswer, studentAnswer string) (EvaluationResult, error)
	AskHelp(ctx context.Context, cfg Config, prompt string) (string, error)
	Proxy(ctx context.Context, cfg Config, method, path, contentType string, body []byte) (*ProxyResponse, error)
}

type gateway struct {
	client *http.Client
}

func NewGateway() Gateway {
	return &gateway{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (g *gateway) IsConfigured(cfg Config) bool {
	return strings.TrimSpace(cfg.BaseURL) != "" && strings.TrimSpace(cfg.APIKey) != ""
}
