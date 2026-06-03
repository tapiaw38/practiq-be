package ai

import (
	"io"
	"net/http"

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
