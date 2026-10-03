package ai

import (
	"github.com/gin-gonic/gin"
	ucAI "github.com/tapiaw38/practiq-be/internal/usecases/ai"
)

func NewProxySendTextMessageHandler(uc ucAI.ProxyUsecase) gin.HandlerFunc {
	return proxyToAssistant(uc, messagePath, textMessageToTutorMultipart)
}
