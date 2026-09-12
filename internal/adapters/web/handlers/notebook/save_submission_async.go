package notebook

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	notebookRepo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/notebook"
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
		// Taken here, where the delivery is accepted, and not inside the
		// goroutine: what has to be ordered is the order the student sent
		// them, and the goroutines finish in whatever order the assistant
		// replies.
		version := now.UnixNano()
		if err := repo.Create(c.Request.Context(), domain.SubmitJob{
			ID:        jobID,
			Kind:      "notebook",
			StudentID: studentID,
			Status:    "processing",
			CreatedAt: now,
			UpdatedAt: now,
		}); err != nil {
			// Returning 202 with a job that was never stored hands the client an
			// id that polls as not-found forever, while the submission still
			// runs — so the student cannot see the result and may retry work
			// that was already applied.
			log.Printf("failed to create submit job: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "practice_sheet:submit-job-not-created",
				"message": "could not start the submission",
			})
			return
		}

		go func(pid, sid, jid string, ver int64, payload struct {
			CanvasData string `json:"canvas_data"`
			AnswerText string `json:"answer_text"`
		}) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			finishCtx, finishCancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer finishCancel()

			err := uc.Execute(ctx, ucNB.SaveSubmissionInput{
				PageID:     pid,
				StudentID:  sid,
				CanvasData: payload.CanvasData,
				AnswerText: payload.AnswerText,
				Version:    ver,
			})
			// A delivery the student has already replaced is not a failure to
			// report. Telling them this one failed would send them to resubmit
			// work that the newer answer already covers.
			if errors.Is(err, notebookRepo.ErrStaleSubmission) {
				if updateErr := repo.Update(finishCtx, domain.SubmitJob{
					ID:     jid,
					Status: "done",
				}); updateErr != nil {
					log.Printf("failed to update submit job: %v", updateErr)
				}
				return
			}
			if err != nil {
				if updateErr := repo.Update(finishCtx, domain.SubmitJob{
					ID:        jid,
					Status:    "failed",
					ErrorCode: "notebook:submit-failed",
					Message:   err.Error(),
				}); updateErr != nil {
					log.Printf("failed to update submit job: %v", updateErr)
				}
				return
			}
			if updateErr := repo.Update(finishCtx, domain.SubmitJob{
				ID:     jid,
				Status: "done",
			}); updateErr != nil {
				log.Printf("failed to update submit job: %v", updateErr)
			}
		}(pageID, studentID, jobID, version, input)

		c.JSON(http.StatusAccepted, gin.H{
			"data": gin.H{
				"job_id": jobID,
				"status": "processing",
			},
		})
	}
}
