package userprofile

import (
	"net/http"

	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
)

func NewFindByEmailHandler(uc ucProfile.FindByEmailUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")
		if email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "email is required"})
			return
		}

		output, appErr := uc.Execute(c, email, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
