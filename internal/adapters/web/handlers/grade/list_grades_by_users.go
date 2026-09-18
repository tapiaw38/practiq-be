package grade

import (
	"net/http"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"

	"github.com/gin-gonic/gin"
)

const maxBatchUserIDs = 200

type listGradesByUsersInput struct {
	UserIDs []string `json:"user_ids"`
}

func NewListGradesByUsersHandler(uc ucGrade.ListGradesByUsersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input listGradesByUsersInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if len(input.UserIDs) > maxBatchUserIDs {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": "too many user_ids"})
			return
		}

		output, appErr := uc.Execute(c, middlewares.GetUserID(c), middlewares.IsTeacher(c), input.UserIDs)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
