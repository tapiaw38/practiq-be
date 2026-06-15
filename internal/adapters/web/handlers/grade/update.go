package grade

import (
	"net/http"

	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

type updateInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewUpdateHandler(uc ucGrade.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucGrade.UpdateInput{
			Name:        input.Name,
			Description: input.Description,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
