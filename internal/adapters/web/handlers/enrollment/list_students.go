package enrollment

import (
	"net/http"
	"strconv"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucEnrollment "github.com/tapiaw38/practiq-be/internal/usecases/enrollment"

	"github.com/gin-gonic/gin"
)

func NewListStudentsHandler(uc ucEnrollment.ListStudentsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")

		input := ucEnrollment.ListStudentsInput{
			CourseID:     courseID,
			RequesterID:  middlewares.GetUserID(c),
			IsSuperAdmin: middlewares.IsSuperAdmin(c),
			BearerToken:  c.GetHeader("Authorization"),
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

		output, appErr := uc.Execute(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
