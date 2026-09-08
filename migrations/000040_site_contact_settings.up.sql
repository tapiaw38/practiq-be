CREATE TABLE site_contact_settings (id SMALLINT PRIMARY KEY CHECK (id = 1), email TEXT NOT NULL, phone TEXT NOT NULL, whatsapp TEXT NOT NULL, updated_at TIMESTAMP NOT NULL DEFAULT NOW());
INSERT INTO site_contact_settings (id,email,phone,whatsapp) VALUES (1,'hola@practiq.com.ar','+54 9 383 743-0889','https://wa.me/5493837430889') ON CONFLICT (id) DO NOTHING;
