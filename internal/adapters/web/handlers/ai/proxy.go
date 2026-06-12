package ai

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func proxyToAssistant(uc ucAI.ProxyUsecase, pathBuilder func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid request body"})
			return
		}
		logAssistantProxyBody(c, body)

		output, appErr := uc.Execute(c, ucAI.ProxyInput{
			UserID:      middlewares.GetUserID(c),
			Method:      c.Request.Method,
			Path:        pathBuilder(c),
			ContentType: c.GetHeader("Content-Type"),
			Body:        body,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		if output.ContentType != "" {
			c.Header("Content-Type", output.ContentType)
		}
		c.Data(output.StatusCode, output.ContentType, output.Body)
	}
}

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

	imageCount := 0
	imageBytes := int64(0)
	imageNames := []string{}
	audioCount := 0
	audioBytes := int64(0)
	audioNames := []string{}
	if req.MultipartForm != nil {
		for field, files := range req.MultipartForm.File {
			for _, file := range files {
				switch field {
				case "image_content":
					imageCount++
					imageBytes += file.Size
					imageNames = append(imageNames, file.Filename)
				case "voice_content":
					audioCount++
					audioBytes += file.Size
					audioNames = append(audioNames, file.Filename)
				default:
					continue
				}
			}
		}
	}

	log.Printf("[assistant_proxy] method=%s path=%s content_type=%q body_bytes=%d image_count=%d image_bytes=%d image_names=%v audio_count=%d audio_bytes=%d audio_names=%v",
		c.Request.Method,
		c.Request.URL.RequestURI(),
		contentType,
		len(body),
		imageCount,
		imageBytes,
		imageNames,
		audioCount,
		audioBytes,
		audioNames,
	)
}

func NewProxyListConversationsHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/user"
	})
}

func NewProxyGetConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/" + c.Param("id")
	})
}

func NewProxyCreateConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		return "/conversation/"
	})
}

func NewProxySendMessageHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string {
		path := "/conversation/" + c.Param("id") + "/message"
		if rawQuery := c.Request.URL.RawQuery; rawQuery != "" {
			path += "?" + rawQuery
		}
		return path
	})
}
