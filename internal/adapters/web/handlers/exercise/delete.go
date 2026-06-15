package exercise

import (
	"net/http"

	ucExercise "github.com/tapiaw38/practiq-be/internal/usecases/exercise"

	"github.com/gin-gonic/gin"
)

func NewDeleteHandler(uc ucExercise.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "exercise deleted"})
	}
}
