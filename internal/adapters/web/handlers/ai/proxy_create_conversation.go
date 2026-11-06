package ai

import (
	"github.com/gin-gonic/gin"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func NewProxyCreateConversationHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, func(*gin.Context) string { return "/conversation/" }, nil)
}
