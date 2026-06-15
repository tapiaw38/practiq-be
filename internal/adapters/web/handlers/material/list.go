package material

import (
	"net/http"

	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucMaterial.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
