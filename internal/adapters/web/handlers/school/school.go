package school

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tapiaw38/practiq-be/internal/adapters/web/middlewares"
	ucSchool "github.com/tapiaw38/practiq-be/internal/usecases/school"
)

type Handler struct{ service ucSchool.Service }

func NewHandler(service ucSchool.Service) *Handler { return &Handler{service: service} }

type createInput struct {
	Name string `json:"name" binding:"required"`
	Slug string `json:"slug" binding:"required"`
}
type memberInput struct {
	UserID         string `json:"user_id" binding:"required"`
	MembershipRole string `json:"membership_role"`
	ProfileType    string `json:"profile_type" binding:"required"`
}

func (h *Handler) Create(c *gin.Context) {
	var in createInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
		return
	}
	out, appErr := h.service.Create(c, ucSchool.CreateInput{Name: in.Name, Slug: in.Slug, CreatedBy: middlewares.GetUserID(c)})
	if appErr != nil {
		appErr.Log(c)
		c.JSON(appErr.StatusCode(), appErr)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *Handler) List(c *gin.Context) {
	out, appErr := h.service.List(c, middlewares.GetUserID(c), middlewares.IsSuperAdmin(c))
	if appErr != nil {
		appErr.Log(c)
		c.JSON(appErr.StatusCode(), appErr)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) Members(c *gin.Context) {
	out, appErr := h.service.Members(c, c.Param("id"), middlewares.GetUserID(c), middlewares.IsSuperAdmin(c))
	if appErr != nil {
		appErr.Log(c)
		c.JSON(appErr.StatusCode(), appErr)
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) AddMember(c *gin.Context) {
	var in memberInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "common:bad-request", "message": err.Error()})
		return
	}
	if in.MembershipRole == "" {
		in.MembershipRole = "member"
	}
	out, appErr := h.service.AddMember(c, ucSchool.AddMemberInput{SchoolID: c.Param("id"), UserID: in.UserID, MembershipRole: in.MembershipRole, ProfileType: in.ProfileType, ActorID: middlewares.GetUserID(c), IsSuperAdmin: middlewares.IsSuperAdmin(c)})
	if appErr != nil {
		appErr.Log(c)
		c.JSON(appErr.StatusCode(), appErr)
		return
	}
	c.JSON(http.StatusCreated, out)
}

func (h *Handler) RemoveMember(c *gin.Context) {
	appErr := h.service.RemoveMember(c, ucSchool.RemoveMemberInput{SchoolID: c.Param("id"), UserID: c.Param("userId"), ActorID: middlewares.GetUserID(c), IsSuperAdmin: middlewares.IsSuperAdmin(c)})
	if appErr != nil {
		appErr.Log(c)
		c.JSON(appErr.StatusCode(), appErr)
		return
	}
	c.Status(http.StatusNoContent)
}
