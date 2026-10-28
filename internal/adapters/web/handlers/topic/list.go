package topic

import (
	"net/http"
	"strconv"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucTopic "github.com/tapiaw38/practiq-be/internal/usecases/topic"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucTopic.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")

		input := ucTopic.ListInput{
			CourseID: courseID,
		}

		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				input.Limit = limit
			}
		}

		if offsetStr := c.Query("offset"); offsetStr != "" {
			if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
				input.Offset = offset
			}
		}

		userID := middlewares.GetUserID(c)
		isSuperAdmin := middlewares.IsSuperAdmin(c)
		output, appErr := uc.Execute(c, userID, isSuperAdmin, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
