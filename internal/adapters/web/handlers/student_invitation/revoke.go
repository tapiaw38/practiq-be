package studentinvitation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucInvitation "github.com/tapiaw38/practiq-be/internal/usecases/student_invitation"
)

func NewRevokeHandler(uc ucInvitation.RevokeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.Execute(c, c.Param("id"), middlewares.GetUserID(c)); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.Status(http.StatusNoContent)
	}
}
