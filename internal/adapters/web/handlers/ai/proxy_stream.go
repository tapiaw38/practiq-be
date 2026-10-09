package ai

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func proxyStreamToAssistant(uc ucAI.ProxyUsecase, pathBuilder func(*gin.Context) string, transform requestTransformer) gin.HandlerFunc {
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
		contentType := c.GetHeader("Content-Type")
		if transform != nil {
			contentType, body, err = transform(contentType, body)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "invalid assistant message"})
				return
			}
		}
		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:streaming-unsupported", "message": "streaming is unavailable"})
			return
		}
		started := false
		status, responseType, appErr := uc.ExecuteStream(c, ucAI.ProxyInput{
			UserID: middlewares.GetUserID(c), Method: c.Request.Method, Path: pathBuilder(c), ContentType: contentType, Body: body,
		}, func(upstreamStatus int, upstreamType string, chunk []byte) error {
			if !started {
				if upstreamType == "" {
					upstreamType = "text/event-stream"
				}
				c.Header("Content-Type", upstreamType)
				c.Header("Cache-Control", "no-cache")
				c.Header("X-Accel-Buffering", "no")
				c.Status(upstreamStatus)
				started = true
			}
			if _, writeErr := c.Writer.Write(chunk); writeErr != nil {
				return writeErr
			}
			flusher.Flush()
			return nil
		})
		if appErr != nil {
			if !started {
				appErr.Log(c)
				c.JSON(appErr.StatusCode(), appErr)
			}
			return
		}
		if !started {
			if responseType != "" {
				c.Header("Content-Type", responseType)
			}
			c.Status(status)
		}
	}
}
