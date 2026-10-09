package attemptreview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucReview "github.com/tapiaw38/practiq-be/internal/usecases/attempt_review"
)

func NewReviewHandler(uc ucReview.ReviewUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucReview.ReviewInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(
			c,
			c.Param("id"),
			middlewares.GetUserID(c),
			middlewares.IsSuperAdmin(c),
			input,
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
