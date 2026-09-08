package sitecontact

import (
	"github.com/gin-gonic/gin"
	repo "github.com/tapiaw38/practiq-be/internal/adapters/datasources/repositories/site_contact"
	"strings"
)

type input struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	WhatsApp string `json:"whatsapp" binding:"required"`
}

func Public(r repo.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		v, e := r.Get(c)
		if e != nil {
			c.JSON(500, gin.H{"message": "site contact unavailable"})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{"email": v.Email, "phone": v.Phone, "whatsapp": v.WhatsApp}})
	}
}
func Get(r repo.Repository) gin.HandlerFunc { return Public(r) }
func Update(r repo.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var in input
		if e := c.ShouldBindJSON(&in); e != nil {
			c.JSON(400, gin.H{"message": e.Error()})
			return
		}
		v, e := r.Save(c, repo.Contact{Email: strings.TrimSpace(in.Email), Phone: strings.TrimSpace(in.Phone), WhatsApp: strings.TrimSpace(in.WhatsApp)})
		if e != nil {
			c.JSON(500, gin.H{"message": "could not save site contact"})
			return
		}
		c.JSON(200, gin.H{"data": gin.H{"email": v.Email, "phone": v.Phone, "whatsapp": v.WhatsApp}})
	}
}
