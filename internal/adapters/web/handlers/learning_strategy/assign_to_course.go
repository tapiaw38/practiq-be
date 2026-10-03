package learningstrategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"
)

func NewAssignToCourseHandler(uc ucLS.AssignToCourseUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID := middlewares.GetUserID(c)
		courseID := c.Param("id")
		isSuperAdmin := middlewares.IsSuperAdmin(c)

		var input ucLS.AssignToCourseInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, courseID, isSuperAdmin, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
