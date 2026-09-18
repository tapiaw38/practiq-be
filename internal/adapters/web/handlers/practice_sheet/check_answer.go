package practicesheet

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	apperrors "github.com/tapiaw38/practiq-be/internal/platform/errors"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
)

func NewCheckAnswerHandler(uc ucPS.CheckAnswerUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucPS.CheckAnswerInput
		if err := c.ShouldBindJSON(&input); err != nil {
			appErr := apperrors.NewBadRequestError("invalid request body")
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		output, appErr := uc.Execute(
			c,
			c.Param("id"),
			c.Param("exerciseId"),
			middlewares.GetUserID(c),
			input,
		)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
