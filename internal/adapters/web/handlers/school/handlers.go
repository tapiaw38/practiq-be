package school

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

func NewListHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.List(c)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

// NewMineHandler serves the school selector: what the asking user belongs to,
// with the role they hold in each.
func NewMineHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.Mine(c, middlewares.GetUserID(c))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewCreateHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSchool.SchoolInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Create(c, input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusCreated, output)
	}
}

func NewUpdateHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSchool.SchoolInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		output, appErr := uc.Update(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), input)
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewListMembersHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		output, appErr := uc.ListMembers(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"))
		if appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.JSON(http.StatusOK, output)
	}
}

func NewAddMemberHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input ucSchool.MemberInput
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		if appErr := uc.AddMember(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), input); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func NewRemoveMemberHandler(uc ucSchool.ManageUsecase) gin.HandlerFunc {
	return func(c *gin.Context) {
		if appErr := uc.RemoveMember(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c), c.Param("id"), c.Param("userId")); appErr != nil {
			appErr.Log(c)
			c.JSON(appErr.StatusCode(), appErr)
			return
		}
		c.Status(http.StatusNoContent)
	}
}
