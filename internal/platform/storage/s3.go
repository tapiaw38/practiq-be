package storage

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/tapiaw38/practiq-be/internal/platform/config"
)

type ImageStorage interface {
	IsConfigured() bool
	UploadDataURI(ctx context.Context, folder, userID, dataURI string) (string, error)
	ResolveDataURI(ctx context.Context, value string) (string, error)
}

type NoopImageStorage struct{}

func (NoopImageStorage) IsConfigured() bool { return false }
func (NoopImageStorage) UploadDataURI(ctx context.Context, folder, userID, dataURI string) (string, error) {
	return strings.TrimSpace(dataURI), nil
}
func (NoopImageStorage) ResolveDataURI(ctx context.Context, value string) (string, error) {
	return strings.TrimSpace(value), nil
}

type S3ImageStorage struct {
	cfg    config.S3Config
	client *http.Client
}

func NewS3ImageStorage(cfg config.S3Config) ImageStorage {
	s := &S3ImageStorage{
		cfg:    cfg,
		client: &http.Client{Timeout: 30 * time.Second},
	}
	if !s.IsConfigured() {
		return NoopImageStorage{}
	}
	return s
}

func (s *S3ImageStorage) IsConfigured() bool {
	return strings.TrimSpace(s.cfg.AWSRegion) != "" &&
		strings.TrimSpace(s.cfg.AWSBucket) != "" &&
		strings.TrimSpace(s.cfg.AWSAccessKeyID) != "" &&
		strings.TrimSpace(s.cfg.AWSSecretAccessKey) != ""
}

func (s *S3ImageStorage) UploadDataURI(ctx context.Context, folder, userID, dataURI string) (string, error) {
	body, contentType, err := decodeDataURI(dataURI)
	if err != nil {
		return strings.TrimSpace(dataURI), err
	}

	key := buildImageKey(folder, userID, extensionForContentType(contentType))
	if err := s.putObject(ctx, key, contentType, body); err != nil {
		return strings.TrimSpace(dataURI), err
	}
	return s.objectURL(key), nil
}

func (s *S3ImageStorage) ResolveDataURI(ctx context.Context, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" || strings.HasPrefix(value, "data:image/") {
		return value, nil
	}

	key, ok := s.keyFromValue(value)
	if !ok {
		return value, nil
	}

	body, contentType, err := s.getObject(ctx, key)
	if err != nil {
		return value, err
	}
	if contentType == "" {
		contentType = http.DetectContentType(body)
	}
	return "data:" + contentType + ";base64," + base64.StdEncoding.EncodeToString(body), nil
}

func decodeDataURI(dataURI string) ([]byte, string, error) {
	dataURI = strings.TrimSpace(dataURI)
	parts := strings.SplitN(dataURI, ",", 2)
	if len(parts) != 2 || !strings.HasPrefix(parts[0], "data:image/") {
		return nil, "", errors.New("invalid image data uri")
	}

	contentType := "image/png"
	meta := strings.TrimPrefix(parts[0], "data:")
	if semi := strings.Index(meta, ";"); semi > 0 {
		contentType = meta[:semi]
	}

	body, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, "", err
	}
	return body, contentType, nil
}

func buildImageKey(folder, userID, ext string) string {
	if ext == "" {
		ext = ".png"
	}
	id := randomHex(16)
	return path.Join("image", cleanPathPart(folder), cleanPathPart(userID), id+ext)
}

func cleanPathPart(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "unknown"
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_':
			b.WriteRune(r)
		default:
			b.WriteByte('_')
		}
	}
	return b.String()
}

func extensionForContentType(contentType string) string {
	switch strings.ToLower(strings.TrimSpace(contentType)) {
	case "image/jpeg", "image/jpg":
		return ".jpg"
	case "image/gif":
		return ".gif"
	case "image/webp":
		return ".webp"
	default:
		return ".png"
	}
}

func randomHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

func (s *S3ImageStorage) putObject(ctx context.Context, key, contentType string, body []byte) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, s.objectURL(key), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	s.sign(req, body)

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("s3 put object status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	return nil
}

func (s *S3ImageStorage) getObject(ctx context.Context, key string) ([]byte, string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.objectURL(key), nil)
	if err != nil {
		return nil, "", err
	}
	s.sign(req, nil)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, "", fmt.Errorf("s3 get object status %d: %s", resp.StatusCode, strings.TrimSpace(string(msg)))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func (s *S3ImageStorage) objectURL(key string) string {
	endpoint := strings.TrimRight(strings.TrimSpace(s.cfg.AWSEndpoint), "/")
	if endpoint != "" {
		return endpoint + "/" + strings.Trim(strings.TrimSpace(s.cfg.AWSBucket), "/") + "/" + escapePath(key)
	}
	bucket := strings.TrimSpace(s.cfg.AWSBucket)
	region := strings.TrimSpace(s.cfg.AWSRegion)
	return "https://" + bucket + ".s3." + region + ".amazonaws.com/" + escapePath(key)
}

func escapePath(key string) string {
	parts := strings.Split(key, "/")
	for i := range parts {
		parts[i] = url.PathEscape(parts[i])
	}
	return strings.Join(parts, "/")
}

func (s *S3ImageStorage) keyFromValue(value string) (string, bool) {
	if strings.HasPrefix(value, "s3://") {
		trimmed := strings.TrimPrefix(value, "s3://")
		prefix := strings.TrimSpace(s.cfg.AWSBucket) + "/"
		if strings.HasPrefix(trimmed, prefix) {
			return strings.TrimPrefix(trimmed, prefix), true
		}
		return "", false
	}
	u, err := url.Parse(value)
	if err != nil || u.Path == "" {
		return "", false
	}
	endpoint := strings.TrimRight(strings.TrimSpace(s.cfg.AWSEndpoint), "/")
	if endpoint != "" {
		eu, err := url.Parse(endpoint)
		if err == nil && strings.EqualFold(u.Host, eu.Host) {
			prefix := "/" + strings.Trim(strings.TrimSpace(s.cfg.AWSBucket), "/") + "/"
			if strings.HasPrefix(u.EscapedPath(), prefix) {
				key, err := url.PathUnescape(strings.TrimPrefix(u.EscapedPath(), prefix))
				return key, err == nil
			}
		}
	}
	bucketHost := strings.TrimSpace(s.cfg.AWSBucket) + ".s3." + strings.TrimSpace(s.cfg.AWSRegion) + ".amazonaws.com"
	if strings.EqualFold(u.Host, bucketHost) {
		key, err := url.PathUnescape(strings.TrimPrefix(u.EscapedPath(), "/"))
		return key, err == nil
	}
	return "", false
}

func (s *S3ImageStorage) sign(req *http.Request, body []byte) {
	now := time.Now().UTC()
	amzDate := now.Format("20060102T150405Z")
	date := now.Format("20060102")
	payloadHash := sha256Hex(body)

	req.Header.Set("X-Amz-Date", amzDate)
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	if token := strings.TrimSpace(s.cfg.AWSSessionToken); token != "" {
		req.Header.Set("X-Amz-Security-Token", token)
	}
	req.Host = req.URL.Host

	signedHeaders, canonicalHeaders := canonicalHeaders(req)
	canonicalRequest := strings.Join([]string{
		req.Method,
		canonicalURI(req.URL),
		canonicalQuery(req.URL),
		canonicalHeaders,
		signedHeaders,
		payloadHash,
	}, "\n")

	scope := date + "/" + strings.TrimSpace(s.cfg.AWSRegion) + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")
	signature := hex.EncodeToString(hmacSHA256(signingKey(s.cfg.AWSSecretAccessKey, date, s.cfg.AWSRegion), []byte(stringToSign)))
	req.Header.Set("Authorization", "AWS4-HMAC-SHA256 Credential="+s.cfg.AWSAccessKeyID+"/"+scope+", SignedHeaders="+signedHeaders+", Signature="+signature)
}

func canonicalHeaders(req *http.Request) (string, string) {
	names := make([]string, 0, len(req.Header)+1)
	values := map[string]string{"host": req.Host}
	for name, vals := range req.Header {
		lower := strings.ToLower(name)
		names = append(names, lower)
		values[lower] = strings.Join(vals, ",")
	}
	names = append(names, "host")
	sort.Strings(names)
	names = unique(names)

	var headers strings.Builder
	for _, name := range names {
		headers.WriteString(name)
		headers.WriteByte(':')
		headers.WriteString(strings.Join(strings.Fields(values[name]), " "))
		headers.WriteByte('\n')
	}
	return strings.Join(names, ";"), headers.String()
}

func unique(values []string) []string {
	out := values[:0]
	for i, value := range values {
		if i == 0 || value != values[i-1] {
			out = append(out, value)
		}
	}
	return out
}

func canonicalURI(u *url.URL) string {
	if u.EscapedPath() == "" {
		return "/"
	}
	return u.EscapedPath()
}

func canonicalQuery(u *url.URL) string {
	q := u.Query()
	if len(q) == 0 {
		return ""
	}
	keys := make([]string, 0, len(q))
	for key := range q {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, key := range keys {
		values := q[key]
		sort.Strings(values)
		for _, value := range values {
			parts = append(parts, url.QueryEscape(key)+"="+url.QueryEscape(value))
		}
	}
	return strings.Join(parts, "&")
}

func signingKey(secret, date, region string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), []byte(date))
	kRegion := hmacSHA256(kDate, []byte(region))
	kService := hmacSHA256(kRegion, []byte("s3"))
	return hmacSHA256(kService, []byte("aws4_request"))
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func sha256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}
