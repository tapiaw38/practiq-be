package grade

import (
	"net/http"

	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

type assignMemberInput struct {
	UserID string `json:"user_id" binding:"required"`
}

func NewAssignMemberHandler(uc ucGrade.AssignMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		var input assignMemberInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, gradeID, input.UserID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
