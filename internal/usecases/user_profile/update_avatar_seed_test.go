package userprofile

import (
	"strings"
	"testing"
)

func TestValidAvatarSeedAcceptsOpaqueTokens(t *testing.T) {
	for _, seed := range []string{
		"robot-7",
		"Bo_ttts9",
		strings.Repeat("a", maxAvatarSeedLength),
		// Empty clears the avatar back to the default.
		"",
	} {
		if !validAvatarSeed(seed) {
			t.Errorf("%q deberia ser una semilla valida", seed)
		}
	}
}

func TestValidAvatarSeedRejectsAnythingResolvable(t *testing.T) {
	// The whole point of storing a seed is that it can never be a URL or
	// markup: whatever lands here is rendered by the client.
	for _, seed := range []string{
		"https://evil.example/pwn.svg",
		"<script>alert(1)</script>",
		"seed with spaces",
		"../../etc/passwd",
		`"onerror="alert(1)`,
		strings.Repeat("a", maxAvatarSeedLength+1),
	} {
		if validAvatarSeed(seed) {
			t.Errorf("%q no deberia aceptarse", seed)
		}
	}
}
