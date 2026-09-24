package secretbox

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

var (
	ErrNoKey      = errors.New("secretbox: no encryption key configured")
	ErrCiphertext = errors.New("secretbox: ciphertext is not valid")
)

// Box seals and opens short secrets with AES-256-GCM.
type Box struct {
	aead cipher.AEAD
}

// New derives the cipher from a passphrase of any length.
func New(key string) (*Box, error) {
	if strings.TrimSpace(key) == "" {
		return nil, ErrNoKey
	}
	sum := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return &Box{aead: aead}, nil
}

// Seal returns base64(nonce || ciphertext).
func (b *Box) Seal(plaintext string) (string, error) {
	nonce := make([]byte, b.aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := b.aead.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(sealed), nil
}

// Open reverses Seal.
func (b *Box) Open(encoded string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", ErrCiphertext
	}
	if len(raw) < b.aead.NonceSize() {
		return "", ErrCiphertext
	}
	nonce, ciphertext := raw[:b.aead.NonceSize()], raw[b.aead.NonceSize():]
	plaintext, err := b.aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", ErrCiphertext
	}
	return string(plaintext), nil
}

// Last4 renders a secret for display without revealing it.
func Last4(secret string) string {
	secret = strings.TrimSpace(secret)
	if len(secret) <= 4 {
		return ""
	}
	return secret[len(secret)-4:]
}
