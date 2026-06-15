package material

import (
	"net/http"

	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"

	"github.com/gin-gonic/gin"
)

func NewDeleteHandler(uc ucMaterial.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "material deleted"})
	}
}
