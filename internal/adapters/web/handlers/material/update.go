package material

import (
	"net/http"

	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
)

type updateMaterialInput struct {
	Title         string `json:"title" binding:"required"`
	ExtractedText string `json:"extracted_text"`
	FileURL       string `json:"file_url"`
}

func NewUpdateHandler(uc ucMaterial.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		requesterID := middlewares.GetUserID(c)
		isAdmin := middlewares.HasRole(c, "admin", "superadmin")

		var input updateMaterialInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, requesterID, isAdmin, id, ucMaterial.UpdateInput{
			Title:         input.Title,
			ExtractedText: input.ExtractedText,
			FileURL:       input.FileURL,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}
