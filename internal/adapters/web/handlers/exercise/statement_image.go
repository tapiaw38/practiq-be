package exercise

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"
)

func NewStatementImageHandler(uc ucExercise.StatementImageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(
			c,
			middlewares.GetUserID(c),
			middlewares.IsSuperAdmin(c),
			c.Param("id"),
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		// Private: the statement belongs to one course. The window is short
		// because a teacher who redraws it expects to see the new one.
		c.Header("Cache-Control", "private, max-age=300")
		c.Data(http.StatusOK, output.ContentType, output.Content)
	}
}
