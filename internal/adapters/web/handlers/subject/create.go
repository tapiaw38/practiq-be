package subject

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSubject "github.com/tapiaw38/practiq-be/internal/usecases/subject"
)

type createInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewCreateHandler(uc ucSubject.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucSubject.CreateInput{
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   middlewares.GetUserID(c),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
