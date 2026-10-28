package ai

import (
	"github.com/gin-gonic/gin"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func NewProxyGetConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(c *gin.Context) string { return "/conversation/" + c.Param("id") }, nil)
}
