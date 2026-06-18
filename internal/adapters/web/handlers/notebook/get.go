package notebook

import (
	"net/http"

	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

func NewGetHandler(uc ucNB.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := c.Query("student_id")
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		out, appErr := uc.Execute(c, requesterID, isAdmin, id, studentID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		if out == nil {
			c.JSON(http.StatusNotFound, gin.H{"message": "notebook not found"})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}
