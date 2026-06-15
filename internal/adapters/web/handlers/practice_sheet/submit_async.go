package practicesheet

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	submitjob "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/submit_job"
	ucPS "github.com/tapiaw38/practiq-be/internal/usecases/practice_sheet"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/tapiaw38/practiq-be/internal/platform/utils"
)

func NewSubmitAsyncHandler(uc ucPS.SubmitUsecase, repo submitjob.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		studentID := middlewares.GetUserID(c)
		var input submitInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		jobID := utils.NewSubmitJobID()
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
