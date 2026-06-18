package material

import (
	"net/http"
	"strconv"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"

	"github.com/gin-gonic/gin"
)

func NewListHandler(uc ucMaterial.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")

		input := ucMaterial.ListInput{
			CourseID: courseID,
		}

		if limitStr := c.Query("limit"); limitStr != "" {
			if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
				input.Limit = limit
			}
		}

		if offsetStr := c.Query("offset"); offsetStr != "" {
			if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
				input.Offset = offset
			}
		}

		userID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		output, appErr := uc.Execute(c, userID, isAdmin, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
