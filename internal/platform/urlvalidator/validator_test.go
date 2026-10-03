package urlvalidator

import (
	"net"
	"strings"
	"testing"
)

func TestValidateIPBlocksPrivateByDefault(t *testing.T) {
	err := validateIP(parseIP(t, "172.17.0.1"), false)
	if err == nil {
		t.Fatal("expected private IP to be blocked")
	}
	if !strings.Contains(err.Error(), "blocked range") {
		t.Fatalf("expected blocked range error, got %v", err)
	}
}

func TestValidateIPAllowsExplicitPrivateHostname(t *testing.T) {
	err := validateIP(parseIP(t, "172.17.0.1"), true)
	if err != nil {
		t.Fatalf("expected private IP to be allowed, got %v", err)
	}
}

func TestMatchesAnyDomain(t *testing.T) {
	if !matchesAnyDomain("host.docker.internal", []string{"host.docker.internal"}) {
		t.Fatal("expected exact private hostname to match")
	}
	if matchesAnyDomain("evil.example.com", []string{"host.docker.internal"}) {
		t.Fatal("expected unrelated hostname not to match")
	}
}

func parseIP(t *testing.T, raw string) net.IP {
	t.Helper()
	ip := net.ParseIP(raw)
	if ip == nil {
		t.Fatalf("invalid test IP %s", raw)
	}
	return ip
}
