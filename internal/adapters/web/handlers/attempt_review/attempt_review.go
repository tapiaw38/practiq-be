package attemptreview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucReview "github.com/tapiaw38/practiq-be/internal/usecases/attempt_review"
)

func NewListHandler(uc ucReview.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c, middlewares.GetUserID(c), c.Query("include_reviewed") == "true")
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
			middlewares.HasRole(c, "admin", "superadmin"),
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
			middlewares.HasRole(c, "admin", "superadmin"),
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
