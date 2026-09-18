package assistant

import (
	"context"
	"net/http"
	"testing"
)

func TestCreateConversationRejectsPrivateAssistantURL(t *testing.T) {
	t.Setenv("ASSISTANT_PROXY_ALLOWED_PRIVATE_HOSTNAMES", "")

	_, err := NewGateway().(*gateway).createConversation(
		context.Background(),
		"http://127.0.0.1",
		"test-key",
	)
	if err == nil {
		t.Fatal("expected private assistant URL to be rejected")
	}
}

func TestGatewayRejectsRedirectToPrivateURL(t *testing.T) {
	t.Setenv("ASSISTANT_PROXY_ALLOWED_PRIVATE_HOSTNAMES", "")

	target, err := http.NewRequest(http.MethodGet, "http://169.254.169.254/latest/meta-data", nil)
	if err != nil {
		t.Fatal(err)
	}
	origin, err := http.NewRequest(http.MethodGet, "https://assistant.example", nil)
	if err != nil {
		t.Fatal(err)
	}

	err = NewGateway().(*gateway).client.CheckRedirect(target, []*http.Request{origin})
	if err == nil {
		t.Fatal("expected private redirect target to be rejected")
	}
}
