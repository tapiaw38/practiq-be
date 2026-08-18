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
		// include_reviewed is kept for the existing clients; `reviewed` is the
		// tri-state the filter bar uses.
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
			CourseID:  c.Query("course_id"),
			StudentID: c.Query("student_id"),
			SheetType: c.Query("sheet_type"),
			Reviewed:  reviewed,
			Limit:     limit,
			Offset:    offset,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewStatementImageHandler(uc ucReview.StatementImageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(
			c,
			c.Param("id"),
			middlewares.GetUserID(c),
			middlewares.IsSuperAdmin(c),
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type reviewInput struct {
	IsCorrect *bool  `json:"is_correct" binding:"required"`
	Feedback  string `json:"feedback"`
}

func NewReviewHandler(uc ucReview.ReviewUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input reviewInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(
			c,
			c.Param("id"),
			middlewares.GetUserID(c),
			middlewares.IsSuperAdmin(c),
			ucReview.ReviewInput{IsCorrect: *input.IsCorrect, Feedback: input.Feedback},
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
