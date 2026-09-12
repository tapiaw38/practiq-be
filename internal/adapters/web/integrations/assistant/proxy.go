package assistant

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// maxAssistantResponseBytes bounds every body read from a user-configured
// assistant endpoint, preventing an unbounded response from exhausting memory.
const maxAssistantResponseBytes = 25 << 20

func (g *gateway) Proxy(ctx context.Context, cfg Config, method, path, contentType string, body []byte) (*ProxyResponse, error) {
	startedAt := time.Now()
	baseURL := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	apiKey := strings.TrimSpace(cfg.APIKey)

	fullURL := baseURL + path
	if err := validateAssistantURL(fullURL); err != nil {
		return nil, fmt.Errorf("URL validation failed: %w", err)
	}

	var requestBody io.Reader
	if len(body) > 0 {
		requestBody = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, requestBody)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := g.client.Do(req)
	if err != nil {
		log.Printf("[assistant_proxy] method=%s path=%s request_bytes=%d duration_ms=%d error=%v", method, path, len(body), time.Since(startedAt).Milliseconds(), err)
		return nil, err
	}
	defer resp.Body.Close()

	responseBody, err := readAssistantResponseBody(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("[assistant_proxy] method=%s path=%s status=%d request_bytes=%d response_bytes=%d duration_ms=%d", method, path, resp.StatusCode, len(body), len(responseBody), time.Since(startedAt).Milliseconds())

	return &ProxyResponse{
		StatusCode:  resp.StatusCode,
		ContentType: resp.Header.Get("Content-Type"),
		Body:        responseBody,
	}, nil
}

func allowedPrivateHostnames() []string {
	raw := strings.TrimSpace(os.Getenv("ASSISTANT_PROXY_ALLOWED_PRIVATE_HOSTNAMES"))
	if raw == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	hostnames := make([]string, 0, len(parts))
	for _, part := range parts {
		hostname := strings.TrimSpace(part)
		if hostname != "" {
			hostnames = append(hostnames, hostname)
		}
	}
	return hostnames
}
