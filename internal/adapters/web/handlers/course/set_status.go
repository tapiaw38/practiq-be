package course

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucCourse "github.com/tapiaw38/practiq-be/internal/usecases/course"
)

type setStatusInput struct {
	Status string `json:"status"`
}

func NewSetStatusHandler(uc ucCourse.SetStatusUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input setStatusInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(
			c,
			middlewares.GetUserID(c),
			middlewares.IsSuperAdmin(c),
			c.Param("id"),
			input.Status,
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
