package exercise

import (
	"net/http"
	"strconv"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucExercise.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		topicID := c.Param("id")

		input := ucExercise.ListInput{
			TopicID: topicID,
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
