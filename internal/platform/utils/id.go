package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

// NewSubmitJobID generates a random hex ID for submit jobs.
func NewSubmitJobID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}
