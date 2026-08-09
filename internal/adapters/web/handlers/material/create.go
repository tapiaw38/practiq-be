package material

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucMaterial "github.com/tapiaw38/practiq-be/internal/usecases/material"
)

type createInput struct {
	Title         string `json:"title" binding:"required"`
	Type          string `json:"type" binding:"required"`
	ExtractedText string `json:"extracted_text"`
	// FileURL comes from POST /uploads; empty means a text-only material.
	FileURL string `json:"file_url"`
}

func NewCreateHandler(uc ucMaterial.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		userID := middlewares.GetUserID(c)
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		isAdmin := middlewares.HasRole(c, "admin", "superadmin")
		output, appErr := uc.Execute(c, userID, isAdmin, ucMaterial.CreateInput{
			CourseID:      courseID,
			TeacherID:     userID,
			Title:         input.Title,
			Type:          input.Type,
			ExtractedText: input.ExtractedText,
			FileURL:       input.FileURL,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}
