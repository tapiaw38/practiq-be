package attemptreview

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucReview "github.com/tapiaw38/practiq-be/internal/usecases/attempt_review"
)

func NewListHandler(uc ucReview.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {

		reviewed := c.Query("reviewed")
		if reviewed == "" && c.Query("include_reviewed") != "true" {
			reviewed = "unreviewed"
		}

		limit, _ := strconv.Atoi(c.Query("limit"))
		offset, _ := strconv.Atoi(c.Query("offset"))
		if offset < 0 {
			offset = 0
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), ucReview.ListInput{
			CourseID:    c.Query("course_id"),
			StudentID:   c.Query("student_id"),
			SheetType:   c.Query("sheet_type"),
			Reviewed:    reviewed,
			Limit:       limit,
			Offset:      offset,
			BearerToken: c.GetHeader("Authorization"),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
