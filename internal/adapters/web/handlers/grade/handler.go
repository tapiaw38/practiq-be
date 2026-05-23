package grade

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucGrade "github.com/tapiaw38/practiq-be/internal/usecases/grade"
)

type createInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewCreateHandler(uc ucGrade.CreateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input createInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, ucGrade.CreateInput{
			Name:        input.Name,
			Description: input.Description,
			CreatedBy:   middlewares.GetUserID(c),
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusCreated, output)
	}
}

func NewListHandler(uc ucGrade.ListUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Execute(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

type updateInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewUpdateHandler(uc ucGrade.UpdateUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input updateInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, id, ucGrade.UpdateInput{
			Name:        input.Name,
			Description: input.Description,
		})
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewDeleteHandler(uc ucGrade.DeleteUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if appErr := uc.Execute(c, id); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "grade deleted"})
	}
}

type assignMemberInput struct {
	UserID string `json:"user_id" binding:"required"`
}

func NewAssignMemberHandler(uc ucGrade.AssignMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		var input assignMemberInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
			return
		}

		output, appErr := uc.Execute(c, gradeID, input.UserID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewListMembersHandler(uc ucGrade.ListMembersUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		output, appErr := uc.Execute(c, gradeID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}

		c.JSON(http.StatusOK, output)
	}
}

func NewRemoveMemberHandler(uc ucGrade.RemoveMemberUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		gradeID := c.Param("id")
		userID := c.Param("userId")
		output, appErr := uc.Execute(c, gradeID, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewListUserGradesHandler(uc ucGrade.ListUserGradesUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("userId")
		output, appErr := uc.Execute(c, userID)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}
