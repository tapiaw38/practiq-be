package material

import "github.com/tapiaw38/practiq-be/internal/domain"

type (
	MaterialData struct {
		ID        string `json:"id"`
		CourseID  string `json:"course_id"`
		TeacherID string `json:"teacher_id"`
		Title     string `json:"title"`
		Type      string `json:"type"`
		FileURL   string `json:"file_url,omitempty"`
		// ViewURL is a short-lived signed URL; the bucket is private, so the
		// canonical FileURL is not openable by a browser on its own.
		ViewURL       string `json:"view_url,omitempty"`
		ExtractedText string `json:"extracted_text,omitempty"`
		// ExtractedTextTruncated says the text above is only the beginning and
		// the whole of it has to be read from the material's own endpoint.
		ExtractedTextTruncated bool   `json:"extracted_text_truncated,omitempty"`
		Status                 string `json:"status"`
		CreatedAt              string `json:"created_at"`
	}
)

func toMaterialData(m domain.Material) MaterialData {
	return MaterialData{
		ID:            m.ID,
		CourseID:      m.CourseID,
		TeacherID:     m.TeacherID,
		Title:         m.Title,
		Type:          m.Type,
		FileURL:       m.FileURL,
		ExtractedText: m.ExtractedText,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// previewChars is how much extracted text a listing carries. The listing clamps
// it to about two lines, so this is already generous; it exists so a short
// material needs no second request to be read in full.
const previewChars = 500

// toMaterialPreview is the listing shape. The repository already cut the text
// to previewChars+1, so a longer value means there is more to fetch.
func toMaterialPreview(m domain.Material) MaterialData {
	data := toMaterialData(m)
	// Cut by runes: a byte-wise cut splits accented characters in half, and the
	// text is Spanish.
	if runes := []rune(data.ExtractedText); len(runes) > previewChars {
		data.ExtractedText = string(runes[:previewChars])
		data.ExtractedTextTruncated = true
	}
	return data
}
