package teacherstudentassignment

import (
	"net/http"
	"strconv"

	ucAssignment "github.com/tapiaw38/practiq-be/internal/usecases/teacher_student_assignment"

	"github.com/gin-gonic/gin"
)

func NewListTeachersHandler(uc ucAssignment.ListTeachersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		studentID := c.Param("studentId")

		input := ucAssignment.ListTeachersInput{
			StudentID:   studentID,
			BearerToken: c.GetHeader("Authorization"),
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
