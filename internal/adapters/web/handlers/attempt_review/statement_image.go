package attemptreview

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucReview "github.com/tapiaw38/practiq-be/internal/usecases/attempt_review"
)

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
