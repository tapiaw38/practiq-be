package ai

import (
	"bytes"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func logAssistantProxyBody(c *gin.Context, body []byte) {
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		log.Printf("[assistant_proxy] method=%s path=%s content_type=%q body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), contentType, len(body))
		return
	}
	req, err := http.NewRequest(c.Request.Method, c.Request.URL.String(), bytes.NewReader(body))
	if err != nil {
		log.Printf("[assistant_proxy] method=%s path=%s multipart_parse_request_error=%v body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), err, len(body))
		return
	}
	req.Header.Set("Content-Type", contentType)
	if err := req.ParseMultipartForm(int64(len(body) + 1024)); err != nil {
		log.Printf("[assistant_proxy] method=%s path=%s multipart_parse_error=%v body_bytes=%d", c.Request.Method, c.Request.URL.RequestURI(), err, len(body))
		return
	}

	imageCount, audioCount, docCount := 0, 0, 0
	imageBytes, audioBytes, docBytes := int64(0), int64(0), int64(0)
	imageNames, audioNames, docNames := []string{}, []string{}, []string{}
	if req.MultipartForm != nil {
		for field, files := range req.MultipartForm.File {
			for _, file := range files {
				switch field {
				case "image_content":
					imageCount, imageBytes, imageNames = imageCount+1, imageBytes+file.Size, append(imageNames, file.Filename)
				case "voice_content":
					audioCount, audioBytes, audioNames = audioCount+1, audioBytes+file.Size, append(audioNames, file.Filename)
				case "document_content":
					docCount, docBytes, docNames = docCount+1, docBytes+file.Size, append(docNames, file.Filename)
				}
			}
		}
	}
	log.Printf("[assistant_proxy] method=%s path=%s content_type=%q body_bytes=%d image_count=%d image_bytes=%d image_names=%v audio_count=%d audio_bytes=%d audio_names=%v doc_count=%d doc_bytes=%d doc_names=%v", c.Request.Method, c.Request.URL.RequestURI(), contentType, len(body), imageCount, imageBytes, imageNames, audioCount, audioBytes, audioNames, docCount, docBytes, docNames)
}
