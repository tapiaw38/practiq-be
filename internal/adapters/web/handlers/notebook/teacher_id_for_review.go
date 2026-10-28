package notebook

import (
	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func teacherIDForReview(c *gin.Context) string {
	if middlewares.IsSuperAdmin(c) {
		return ""
	}
	return middlewares.GetUserID(c)
}
