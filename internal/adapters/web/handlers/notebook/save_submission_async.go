package notebook

import (
	"context"
	"log"
	"net/http"
	"time"

	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	ucNB "github.com/tapiaw38/practiq-be/internal/usecases/notebook"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

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

		jobID := utils.NewSubmitJobID()
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
