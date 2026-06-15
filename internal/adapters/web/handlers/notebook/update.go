package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
)

func NewUpdateHandler(uc ucNB.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		out, err := uc.Execute(c, id, ucNB.UpdateInput{
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}
