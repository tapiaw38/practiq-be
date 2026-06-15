package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
)

func NewGetHandler(uc ucNB.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := c.Query("student_id")
		out, err := uc.Execute(c, id, studentID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		if out == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "notebook not found"})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}
