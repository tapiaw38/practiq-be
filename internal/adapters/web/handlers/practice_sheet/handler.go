package practicesheet

import (
	"context"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"
)

type submitJob struct {
	Status    string             `json:"status"`
	Result    *ucPS.SubmitOutput `json:"result,omitempty"`
	ErrorCode string             `json:"error_code,omitempty"`
	Message   string             `json:"message,omitempty"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
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

func NewSubmitAsyncHandler(uc ucPS.SubmitUsecase) gin.HandlerFunc {
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
		setSubmitJob(jobID, submitJob{
			Status:    "processing",
			CreatedAt: now,
			UpdatedAt: now,
		})

		go func(sheetID, uid, jid string, payload submitInput) {
			output, appErr := uc.Execute(context.Background(), sheetID, uid, ucPS.SubmitInput{Attempts: payload.Attempts})
			updated := time.Now().UTC()
			if appErr != nil {
				setSubmitJob(jid, submitJob{
					Status:    "failed",
					ErrorCode: "practice_sheet:submit-failed",
					Message:   appErr.Error(),
					CreatedAt: now,
					UpdatedAt: updated,
				})
				return
			}
			setSubmitJob(jid, submitJob{
				Status:    "done",
				Result:    output,
				CreatedAt: now,
				UpdatedAt: updated,
			})
		}(id, studentID, jobID, input)

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
			c.JSON(http.StatusNotFound, gin.H{"code": "practice_sheet:submit-job-not-found", "message": "submit job not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": job})
	}
}
