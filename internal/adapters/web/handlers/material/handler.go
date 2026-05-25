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

		output, appErr := uc.Execute(c, ucMaterial.CreateInput{
			CourseID:      courseID,
			TeacherID:     userID,
			Title:         input.Title,
			Type:          input.Type,
			ExtractedText: input.ExtractedText,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

func NewListHandler(uc ucMaterial.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		output, appErr := uc.Execute(c, courseID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type updateMaterialInput struct {
	Title         string `json:"title" binding:"required"`
	ExtractedText string `json:"extracted_text"`
}

func NewUpdateHandler(uc ucMaterial.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateMaterialInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucMaterial.UpdateInput{
			Title:         input.Title,
			ExtractedText: input.ExtractedText,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucMaterial.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "material deleted"})
	}
}
