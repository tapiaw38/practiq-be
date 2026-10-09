package userprofile

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucProfile "github.com/tapiaw38/practiq-be/internal/usecases/user_profile"

	"github.com/gin-gonic/gin"
)

func NewGetByIDHandler(uc ucProfile.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		profileID := c.Param("id")

		if profileID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "profile id required"})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, profileID, c.GetHeader("Authorization"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
