package assistant

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/urlvalidator"
)

type Gateway interface {
	IsConfigured(cfg Config) bool
	AnalyzeCanvas(ctx context.Context, cfg Config, canvasData, correctAnswer string) (string, error)
	AnalyzeNotebookCanvas(ctx context.Context, cfg Config, canvasData, pageContext string) (string, error)
	AnalyzeNotebookStatement(ctx context.Context, cfg Config, imageData, pageContext string) (string, error)
	EvaluatePracticeAnswer(ctx context.Context, cfg Config, question, correctAnswer, studentAnswer, gradeName string) (EvaluationResult, error)
	EvaluateAttachment(ctx context.Context, cfg Config, input AttachmentEvaluationInput) (EvaluationResult, error)
	AskHelp(ctx context.Context, cfg Config, prompt string) (string, error)
	GenerateCourseCuriosities(ctx context.Context, cfg Config, subject, topic, gradeName string, count int) ([]string, error)
	Proxy(ctx context.Context, cfg Config, method, path, contentType string, body []byte) (*ProxyResponse, error)
}

type gateway struct {
	client *http.Client
}

// maxAssistantRedirects is generous for a legitimate endpoint and short enough
// that a redirect loop cannot tie up a request for the full timeout.
const maxAssistantRedirects = 5

func NewGateway() Gateway {
	return &gateway{
		client: &http.Client{
			Timeout: 5 * time.Minute,
			// Every URL this client reaches is user-configured, and validating
			// only the first one left the door open: an attacker-controlled
			// assistant could answer with a redirect to a link-local address
			// like 169.254.169.254, and the default client would follow it and
			// hand the internal response back to the caller. Each hop is
			// revalidated with the same rules as the original.
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= maxAssistantRedirects {
					return fmt.Errorf("stopped after %d redirects", maxAssistantRedirects)
				}
				if err := urlvalidator.ValidateURLWithOptions(req.URL.String(), urlvalidator.Options{
					AllowedPrivateHostnames: allowedPrivateHostnames(),
				}); err != nil {
					return fmt.Errorf("redirect target rejected: %w", err)
				}
				return nil
			},
		},
	}
}

func (g *gateway) IsConfigured(cfg Config) bool {
	return strings.TrimSpace(cfg.BaseURL) != "" && strings.TrimSpace(cfg.APIKey) != ""
}
