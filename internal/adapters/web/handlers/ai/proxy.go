package ai

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

type requestTransformer func(contentType string, body []byte) (string, []byte, error)

const maxAssistantProxyRequestBytes = 25 << 20

func proxyToAssistant(uc ucAI.ProxyUsecase, pathBuilder func(*gin.Context) string, transform requestTransformer) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxAssistantProxyRequestBytes)
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			if _, ok := err.(*http.MaxBytesError); ok {
				c.JSON(http.StatusRequestEntityTooLarge, gin.H{"code": "common:payload-too-large", "message": "assistant attachment exceeds 25 MiB"})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid request body"})
			return
		}
		logAssistantProxyBody(c, body)
		contentType := c.GetHeader("Content-Type")
		if transform != nil {
			contentType, body, err = transform(contentType, body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid assistant message"})
				return
			}
		}

		output, appErr := uc.Execute(c, ucAI.ProxyInput{
			UserID:      middlewares.GetUserID(c),
			Method:      c.Request.Method,
			Path:        pathBuilder(c),
			ContentType: contentType,
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
