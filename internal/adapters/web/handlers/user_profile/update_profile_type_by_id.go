package userprofile

import (
	"net/http"

	"github.com/gin-gonic/gin"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"
)

type profileTypeInput struct {
	ProfileType string `json:"profile_type" binding:"required"`
}

func NewUpdateProfileTypeByIDHandler(uc ucProfile.UpdateProfileTypeUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input profileTypeInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		output, appErr := uc.Execute(c, c.Param("id"), input.ProfileType, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
