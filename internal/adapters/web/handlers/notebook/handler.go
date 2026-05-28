package notebook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

type submitJob struct {
	Status    string    `json:"status"`
	ErrorCode string    `json:"error_code,omitempty"`
	Message   string    `json:"message,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var submitJobs = struct {
	mu   sync.RWMutex
	data map[string]submitJob
}{
	data: make(map[string]submitJob),
}

func newSubmitJobID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

func setSubmitJob(id string, job submitJob) {
	submitJobs.mu.Lock()
	defer submitJobs.mu.Unlock()
	submitJobs.data[id] = job
}

func getSubmitJob(id string) (submitJob, bool) {
	submitJobs.mu.RLock()
	defer submitJobs.mu.RUnlock()
	job, ok := submitJobs.data[id]
	return job, ok
}

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

func NewUpdateHandler(uc ucNB.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input struct {
			Title       string `json:"title" binding:"required"`
			Description string `json:"description"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		out, err := uc.Execute(c, id, ucNB.UpdateInput{
			Title:       input.Title,
			Description: input.Description,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, out)
	}
}

func NewDeleteHandler(uc ucNB.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := uc.Execute(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "notebook deleted"})
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

func NewSaveSubmissionAsyncHandler(uc ucNB.SaveSubmissionUsecase) gin.HandlerFunc {
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

		jobID := newSubmitJobID()
		now := time.Now().UTC()
		setSubmitJob(jobID, submitJob{Status: "processing", CreatedAt: now, UpdatedAt: now})

		go func(pid, sid, jid string, payload struct {
			CanvasData string `json:"canvas_data"`
			AnswerText string `json:"answer_text"`
		}) {
			err := uc.Execute(context.Background(), ucNB.SaveSubmissionInput{
				PageID:     pid,
				StudentID:  sid,
				CanvasData: payload.CanvasData,
				AnswerText: payload.AnswerText,
			})
			updated := time.Now().UTC()
			if err != nil {
				setSubmitJob(jid, submitJob{
					Status:    "failed",
					ErrorCode: "notebook:submit-failed",
					Message:   err.Error(),
					CreatedAt: now,
					UpdatedAt: updated,
				})
				return
			}
			setSubmitJob(jid, submitJob{Status: "done", CreatedAt: now, UpdatedAt: updated})
		}(pageID, studentID, jobID, input)

		c.JSON(http.StatusAccepted, gin.H{
			"data": gin.H{
				"job_id": jobID,
				"status": "processing",
			},
		})
	}
}

func NewGetSubmitJobHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("jobId")
		job, ok := getSubmitJob(jobID)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submit-job-not-found", "message": "submit job not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": job})
	}
}
