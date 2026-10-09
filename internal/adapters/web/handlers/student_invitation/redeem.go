package studentinvitation

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucInvitation "github.com/tapiaw38/practiq-be/internal/usecases/student_invitation"
)

func NewRedeemHandler(uc ucInvitation.RedeemUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucInvitation.RedeemInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.GetHeader("Authorization"), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
