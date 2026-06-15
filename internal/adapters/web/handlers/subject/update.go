package subject

import (
	"net/http"

	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"

	"github.com/gin-gonic/gin"
)

type updateInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewUpdateHandler(uc ucSubject.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucSubject.UpdateInput{
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
