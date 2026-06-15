package learningstrategy

import (
	"net/http"

	ucLS "github.com/tapiaw38/practiq-be/internal/usecases/learning_strategy"

	"github.com/gin-gonic/gin"
)

type updateInput struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func NewUpdateHandler(uc ucLS.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucLS.UpdateInput{
			Name:        input.Name,
			Code:        input.Code,
			Description: input.Description,
			Status:      input.Status,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
