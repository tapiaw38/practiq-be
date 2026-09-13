package sitecontact

import (
	"context"
	"database/sql"
)

type Contact struct {
	Email, Phone, WhatsApp string `json:"-"`
}
type Repository interface {
	Get(context.Context) (Contact, error)
	Save(context.Context, Contact) (Contact, error)
}
type repository struct{ db *sql.DB }

func NewRepository(db *sql.DB) Repository { return &repository{db} }
func (r *repository) Get(ctx context.Context) (v Contact, err error) {
	err = r.db.QueryRowContext(ctx, `SELECT email,phone,whatsapp FROM site_contact_settings WHERE id=1`).Scan(&v.Email, &v.Phone, &v.WhatsApp)
	return
}
func (r *repository) Save(ctx context.Context, v Contact) (Contact, error) {
	_, err := r.db.ExecContext(ctx, `INSERT INTO site_contact_settings(id,email,phone,whatsapp,updated_at) VALUES(1,$1,$2,$3,NOW()) ON CONFLICT(id) DO UPDATE SET email=EXCLUDED.email,phone=EXCLUDED.phone,whatsapp=EXCLUDED.whatsapp,updated_at=NOW()`, v.Email, v.Phone, v.WhatsApp)
	return v, err
}
