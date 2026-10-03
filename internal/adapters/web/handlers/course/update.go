package course

import (
	"net/http"

	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewUpdateHandler(uc ucCourse.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		var input ucCourse.UpdateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isSuperAdmin, id, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
