package notebook

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

func NewCreateHandler(uc ucNB.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		teacherID := middlewares.GetUserID(c)
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
			Level       int    `json:"level"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		out, err := uc.Execute(c, teacherID, ucNB.CreateInput{
			CourseID:    courseID,
			Title:       input.Title,
			Description: input.Description,
			Level:       input.Level,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}

func NewListHandler(uc ucNB.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		out, err := uc.Execute(c, courseID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

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

func NewAddPageHandler(uc ucNB.AddPageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		notebookID := c.Param("id")
		var input struct {
			PageNumber   int    `json:"page_number"`
			Title        string `json:"title"`
			ContentType  string `json:"content_type"`
			ContentData  string `json:"content_data"`
			Instructions string `json:"instructions"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		out, err := uc.Execute(c, ucNB.AddPageInput{
			NotebookID:   notebookID,
			PageNumber:   input.PageNumber,
			Title:        input.Title,
			ContentType:  input.ContentType,
			ContentData:  input.ContentData,
			Instructions: input.Instructions,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, out)
	}
}

func NewUpdatePageHandler(uc ucNB.UpdatePageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		var input struct {
			Title        string `json:"title"`
			ContentType  string `json:"content_type"`
			ContentData  string `json:"content_data"`
			Instructions string `json:"instructions"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if err := uc.Execute(c, ucNB.UpdatePageInput{
			PageID:       pageID,
			Title:        input.Title,
			ContentType:  input.ContentType,
			ContentData:  input.ContentData,
			Instructions: input.Instructions,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}

func NewSaveSubmissionHandler(uc ucNB.SaveSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		pageID := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input struct {
			CanvasData string `json:"canvas_data"`
			AnswerText string `json:"answer_text"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		if err := uc.Execute(c, ucNB.SaveSubmissionInput{
			PageID:     pageID,
			StudentID:  studentID,
			CanvasData: input.CanvasData,
			AnswerText: input.AnswerText,
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, nil)
	}
}
