package sitecontact

import "github.com/tapiaw38/practiq-be/internal/domain"

type ContactData struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	WhatsApp string `json:"whatsapp"`
}

func toContactData(contact domain.SiteContact) ContactData {
	return ContactData{Email: contact.Email, Phone: contact.Phone, WhatsApp: contact.WhatsApp}
}
