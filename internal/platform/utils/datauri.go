package utils

import (
	"encoding/base64"
	"errors"
	"strings"
)

// DecodeDataURI decodes a base64 data URI (e.g., "data:image/png;base64,...")
// Returns the decoded bytes and content type.
func DecodeDataURI(dataURI string) ([]byte, string, error) {
	dataURI = strings.TrimSpace(dataURI)
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 {
		return nil, "", errors.New("invalid data uri: missing comma separator")
	}

	meta := parts[0]
	payload := parts[1]
	contentType := "image/png"

	if strings.HasPrefix(meta, "data:") {
		typeEnd := strings.Index(meta, ";")
		if typeEnd > len("data:") {
			contentType = meta[len("data:"):typeEnd]
		}
	}

	decoded, err := base64.StdEncoding.DecodeString(payload)
	if err != nil {
		return nil, "", err
	}

	return decoded, contentType, nil
}
