package studentreport

import (
	"bytes"
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/tapiaw38/practiq-be/internal/domain"
)

const (
	marginLeft   = 15.0
	marginRight  = 15.0
	marginTop    = 15.0
	pageWidth    = 210.0
	contentWidth = pageWidth - marginLeft - marginRight
)

type PDFBuilder struct {
	pdf  *gofpdf.Fpdf
	data *domain.StudentReportData
}

func NewPDFBuilder(data *domain.StudentReportData) *PDFBuilder {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(true, 20)
	return &PDFBuilder{pdf: pdf, data: data}
}

func (b *PDFBuilder) Build() ([]byte, error) {
	b.pdf.AddPage()
	b.renderHeader()
	b.renderSummary()
	b.renderMasteryChart()
	b.renderCourseTable()
	b.renderDailyChart()
	b.renderAttemptsTable()

	var buf bytes.Buffer
	if err := b.pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *PDFBuilder) renderHeader() {
	pdf := b.pdf
	data := b.data

	// Title
	pdf.SetFont("Arial", "B", 18)
	pdf.SetTextColor(124, 58, 237) // violet
	pdf.CellFormat(contentWidth, 10, "Reporte de Progreso", "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// Student info
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(0, 0, 0)
	pdf.CellFormat(contentWidth, 6, data.Student.Name, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(contentWidth, 5, data.Student.Email, "", 1, "C", false, 0, "")
	pdf.Ln(2)

	// Period and generation date
	pdf.SetFont("Arial", "", 9)
	periodStr := "Todo el periodo"
	if data.Period.From != nil && data.Period.To != nil {
		periodStr = fmt.Sprintf("%s - %s", data.Period.From.Format("02/01/2006"), data.Period.To.Format("02/01/2006"))
	} else if data.Period.From != nil {
		periodStr = fmt.Sprintf("Desde %s", data.Period.From.Format("02/01/2006"))
	} else if data.Period.To != nil {
		periodStr = fmt.Sprintf("Hasta %s", data.Period.To.Format("02/01/2006"))
	}
	pdf.CellFormat(contentWidth, 5, fmt.Sprintf("Periodo: %s | Generado: %s", periodStr, data.GeneratedAt.Format("02/01/2006 15:04")), "", 1, "C", false, 0, "")
	pdf.Ln(6)

	// Divider
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(marginLeft, pdf.GetY(), pageWidth-marginRight, pdf.GetY())
	pdf.Ln(4)
}

func (b *PDFBuilder) renderSummary() {
	pdf := b.pdf
	s := b.data.Summary

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(contentWidth, 7, "Resumen", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Summary table (3x2)
	colWidth := contentWidth / 3
	rowHeight := 12.0

	cells := [][]string{
		{"Temas practicados", fmt.Sprintf("%d", s.TopicsPracticed)},
		{"Dominio promedio", fmt.Sprintf("%.0f%%", s.AverageMastery)},
		{"Intentos totales", fmt.Sprintf("%d", s.TotalAttempts)},
		{"Respuestas correctas", fmt.Sprintf("%d", s.CorrectAttempts)},
		{"Tasa de acierto", fmt.Sprintf("%.0f%%", s.AccuracyRate)},
		{"Racha actual", fmt.Sprintf("%d dias", s.CurrentStreak)},
	}

	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			idx := i*3 + j
			cell := cells[idx]

			pdf.SetFillColor(248, 250, 252) // light gray bg
			pdf.SetFont("Arial", "", 8)
			pdf.SetTextColor(100, 100, 100)

			x := marginLeft + float64(j)*colWidth
			y := pdf.GetY()

			pdf.Rect(x, y, colWidth-2, rowHeight, "F")
			pdf.SetXY(x+2, y+1)
			pdf.CellFormat(colWidth-4, 4, cell[0], "", 0, "L", false, 0, "")

			pdf.SetFont("Arial", "B", 11)
			pdf.SetTextColor(30, 30, 30)
			pdf.SetXY(x+2, y+5)
			pdf.CellFormat(colWidth-4, 6, cell[1], "", 0, "L", false, 0, "")
		}
		pdf.Ln(rowHeight + 2)
	}
	pdf.Ln(4)
}

func (b *PDFBuilder) renderMasteryChart() {
	if len(b.data.TopicProgress) == 0 {
		return
	}

	chartData, err := GenerateMasteryBarChart(b.data.TopicProgress, 500, 250)
	if err != nil || chartData == nil {
		return
	}

	pdf := b.pdf
	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(contentWidth, 7, "Dominio por Tema", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	imgName := "mastery_chart"
	opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader(imgName, opt, bytes.NewReader(chartData))
	pdf.ImageOptions(imgName, marginLeft, pdf.GetY(), contentWidth, 0, false, opt, 0, "")

	pdf.Ln(65)
}

func (b *PDFBuilder) renderCourseTable() {
	if len(b.data.CourseProgress) == 0 {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(50)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(contentWidth, 7, "Progreso por Curso", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table header
	colWidths := []float64{60, 25, 30, 35, 30}
	headers := []string{"Curso", "Nivel", "Temas", "Dominio", "Ultima act."}

	pdf.SetFillColor(124, 58, 237)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 9)

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 7, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetTextColor(50, 50, 50)
	pdf.SetFont("Arial", "", 9)

	for i, c := range b.data.CourseProgress {
		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		title := c.CourseTitle
		if len(title) > 25 {
			title = title[:22] + "..."
		}

		lastAct := "-"
		if c.LastActivity != nil {
			lastAct = c.LastActivity.Format("02/01/2006")
		}

		cells := []string{
			title,
			fmt.Sprintf("%d", c.CurrentLevel),
			fmt.Sprintf("%d", c.TopicCount),
			fmt.Sprintf("%.0f%%", c.AverageMastery),
			lastAct,
		}

		for j, cell := range cells {
			pdf.CellFormat(colWidths[j], 6, cell, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(6)
}

func (b *PDFBuilder) renderDailyChart() {
	if len(b.data.DailyAttempts) == 0 {
		return
	}

	chartData, err := GenerateDailyAttemptsLineChart(b.data.DailyAttempts, 500, 200)
	if err != nil || chartData == nil {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(70)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(contentWidth, 7, "Actividad Diaria", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	imgName := "daily_chart"
	opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader(imgName, opt, bytes.NewReader(chartData))
	pdf.ImageOptions(imgName, marginLeft, pdf.GetY(), contentWidth, 0, false, opt, 0, "")

	pdf.Ln(55)
}

func (b *PDFBuilder) renderAttemptsTable() {
	if len(b.data.RecentAttempts) == 0 {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(50)

	pdf.SetFont("Arial", "B", 11)
	pdf.SetTextColor(50, 50, 50)
	pdf.CellFormat(contentWidth, 7, "Ultimos Intentos (max 50)", "", 1, "L", false, 0, "")
	pdf.Ln(2)

	// Table header
	colWidths := []float64{35, 45, 25, 35, 40}
	headers := []string{"Fecha", "Ejercicio", "Correcto", "Puntaje", "Tiempo"}

	pdf.SetFillColor(124, 58, 237)
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 8)

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 6, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 8)

	attempts := b.data.RecentAttempts
	if len(attempts) > 50 {
		attempts = attempts[:50]
	}

	for i, a := range attempts {
		b.checkPageBreak(8)

		if i%2 == 0 {
			pdf.SetFillColor(248, 250, 252)
		} else {
			pdf.SetFillColor(255, 255, 255)
		}

		correct := "No"
		if a.IsCorrect {
			correct = "Si"
			pdf.SetTextColor(34, 197, 94)
		} else {
			pdf.SetTextColor(239, 68, 68)
		}

		exID := a.ExerciseID
		if len(exID) > 8 {
			exID = exID[:8]
		}

		cells := []string{
			a.CreatedAt.Format("02/01 15:04"),
			exID,
			correct,
			fmt.Sprintf("%.0f%%", a.Score*100),
			fmt.Sprintf("%ds", a.TimeSpentSecs),
		}

		for j, cell := range cells {
			if j == 2 {
				// Keep color for correct/incorrect
			} else {
				pdf.SetTextColor(50, 50, 50)
			}
			pdf.CellFormat(colWidths[j], 5, cell, "1", 0, "C", true, 0, "")
		}
		pdf.Ln(-1)
	}
}

func (b *PDFBuilder) checkPageBreak(height float64) {
	if b.pdf.GetY()+height > 280 {
		b.pdf.AddPage()
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
}
