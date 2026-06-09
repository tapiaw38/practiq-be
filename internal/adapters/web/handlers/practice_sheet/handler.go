package practicesheet

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
)

func newSubmitJobID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

type createInput struct {
	TopicID     string   `json:"topic_id"`
	StrategyID  string   `json:"strategy_id"`
	Title       string   `json:"title" binding:"required"`
	Level       int      `json:"level"`
	SheetType   string   `json:"sheet_type"`
	TestStyle   string   `json:"test_style"`
	ExerciseIDs []string `json:"exercise_ids"`
}

func NewCreateHandler(uc ucPS.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		courseID := c.Param("id")
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucPS.CreateInput{
			CourseID:    courseID,
			SheetType:   input.SheetType,
			TestStyle:   input.TestStyle,
			TopicID:     input.TopicID,
			StrategyID:  input.StrategyID,
			Title:       input.Title,
			Level:       input.Level,
			ExerciseIDs: input.ExerciseIDs,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

func NewListHandler(uc ucPS.ListUsecase) gin.HandlerFunc {
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

func NewGetHandler(uc ucPS.GetUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		output, appErr := uc.Execute(c, id)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type updateSheetInput struct {
	Title       string   `json:"title" binding:"required"`
	TopicID     string   `json:"topic_id"`
	Level       int      `json:"level"`
	SheetType   string   `json:"sheet_type"`
	TestStyle   string   `json:"test_style"`
	ExerciseIDs []string `json:"exercise_ids"`
}

func NewUpdateHandler(uc ucPS.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateSheetInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucPS.UpdateInput{
			Title:       input.Title,
			TopicID:     input.TopicID,
			Level:       input.Level,
			SheetType:   input.SheetType,
			TestStyle:   input.TestStyle,
			ExerciseIDs: input.ExerciseIDs,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucPS.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "practice sheet deleted"})
	}
}

type submitInput struct {
	Attempts []ucPS.AttemptInput `json:"attempts"`
}

func NewSubmitHandler(uc ucPS.SubmitUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input submitInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, studentID, ucPS.SubmitInput{
			Attempts: input.Attempts,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewSubmitAsyncHandler(uc ucPS.SubmitUsecase, repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input submitInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		jobID := newSubmitJobID()
		now := time.Now().UTC()
		if err := repo.Create(c.Request.Context(), domain.SubmitJob{
			ID:        jobID,
			Kind:      "practice_sheet",
			Status:    "processing",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			log.Printf("failed to create submit job: %v", err)
		}

		go func(sheetID, uid, jid string, payload submitInput) {
			output, appErr := uc.Execute(context.Background(), sheetID, uid, ucPS.SubmitInput{Attempts: payload.Attempts})
			if appErr != nil {
				if err := repo.Update(context.Background(), domain.SubmitJob{
					ID:        jid,
					Status:    "failed",
					ErrorCode: "practice_sheet:submit-failed",
					Message:   appErr.Error(),
				}); err != nil {
					log.Printf("failed to update submit job: %v", err)
				}
				return
			}
			resultJSON, _ := json.Marshal(output)
			if err := repo.Update(context.Background(), domain.SubmitJob{
				ID:     jid,
				Status: "done",
				Result: resultJSON,
			}); err != nil {
				log.Printf("failed to update submit job: %v", err)
			}
		}(id, studentID, jobID, input)

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
			c.JSON(http.StatusNotFound, gin.H{"code": "practice_sheet:submit-job-not-found", "message": "submit job not found"})
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
		if len(job.Result) > 0 {
			var result ucPS.SubmitOutput
			if err := json.Unmarshal(job.Result, &result); err == nil {
				response["result"] = &result
			}
		}
		c.JSON(http.StatusOK, gin.H{"data": response})
	}
}
