package notebook

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"
)

func newSubmitJobID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
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

func NewSaveSubmissionAsyncHandler(uc ucNB.SaveSubmissionUsecase, repo submitjob.Repository) gin.HandlerFunc {
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
		if err := repo.Create(c.Request.Context(), domain.SubmitJob{
			ID:        jobID,
			Kind:      "notebook",
			Status:    "processing",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			log.Printf("failed to create submit job: %v", err)
		}

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
			if err != nil {
				if updateErr := repo.Update(context.Background(), domain.SubmitJob{
					ID:        jid,
					Status:    "failed",
					ErrorCode: "notebook:submit-failed",
					Message:   err.Error(),
				}); updateErr != nil {
					log.Printf("failed to update submit job: %v", updateErr)
				}
				return
			}
			if updateErr := repo.Update(context.Background(), domain.SubmitJob{
				ID:     jid,
				Status: "done",
			}); updateErr != nil {
				log.Printf("failed to update submit job: %v", updateErr)
			}
		}(pageID, studentID, jobID, input)

		c.JSON(http.StatusAccepted, gin.H{
			"data": gin.H{
				"job_id": jobID,
				"status": "processing",
			},
		})
	}
}

func NewGetSubmitJobHandler(repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		jobID := c.Param("jobId")
		job, err := repo.GetByID(c.Request.Context(), jobID)
		if err != nil {
			log.Printf("failed to get submit job: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"code": "common:internal-error", "message": "failed to get submit job"})
			return
		}
		if job == nil {
			c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submit-job-not-found", "message": "submit job not found"})
			return
		}
		// Preserve original JSON response shape
		response := gin.H{
			"status":     job.Status,
			"created_at": job.CreatedAt,
			"updated_at": job.UpdatedAt,
		}
		if job.ErrorCode != "" {
			response["error_code"] = job.ErrorCode
		}
		if job.Message != "" {
			response["message"] = job.Message
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}

func NewListSubmissionsHandler(uc ucNB.ListSubmissionsUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		reviewed := ""
		if raw, ok := c.GetQuery("reviewed"); ok {
			if raw == "true" {
				reviewed = "reviewed"
			} else if raw == "false" {
				reviewed = "unreviewed"
			}
		}
		teacherID := ""
		if !middlewares.HasRole(c, "admin", "superadmin") {
			teacherID = middlewares.GetUserID(c)
		}
		output, err := uc.Execute(c, ucNB.ListSubmissionsInput{
			NotebookID: strings.TrimSpace(c.Query("notebook_id")),
			StudentID:  strings.TrimSpace(c.Query("student_id")),
			CourseID:   strings.TrimSpace(c.Query("course_id")),
			Reviewed:   reviewed,
			TeacherID:  teacherID,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:list-submissions-error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewReviewSubmissionHandler(uc ucNB.ReviewSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		teacherID := ""
		if !middlewares.HasRole(c, "admin", "superadmin") {
			teacherID = middlewares.GetUserID(c)
		}
		output, err := uc.Execute(c, submissionID, teacherID)
		if err != nil {
			if err.Error() == "submission not found" {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submission-not-found", "message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:review-error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": output})
	}
}

func NewTeacherReviewSubmissionHandler(uc ucNB.TeacherReviewSubmissionUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		submissionID := c.Param("id")
		var input struct {
			TeacherIsCorrect bool   `json:"teacher_is_correct"`
			TeacherFeedback  string `json:"teacher_feedback"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}
		output, err := uc.Execute(c, submissionID, ucNB.TeacherReviewInput{
			IsCorrect: input.TeacherIsCorrect,
			Feedback:  input.TeacherFeedback,
			TeacherID: teacherIDForReview(c),
		})
		if err != nil {
			if err.Error() == "submission not found" {
				c.JSON(http.StatusNotFound, gin.H{"code": "notebook:submission-not-found", "message": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"code": "notebook:teacher-review-error", "message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": output})
	}
}

func teacherIDForReview(c *gin.Context) string {
	if middlewares.HasRole(c, "admin", "superadmin") {
		return ""
	}
	return middlewares.GetUserID(c)
}
