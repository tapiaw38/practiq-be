package userprofile

import (
	"net/http"

	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
)

func NewGetByIDHandler(uc ucProfile.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		profileID := c.Param("id")
		if profileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "profile id required"})
			return
		}

		output, appErr := uc.Execute(c, profileID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
