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
	marginTop    = 20.0
	marginBottom = 20.0
	pageWidth    = 210.0
	pageHeight   = 297.0
	contentWidth = pageWidth - marginLeft - marginRight
)

// Colors
var (
	colorPrimary   = [3]int{99, 102, 241}  // Indigo
	colorSecondary = [3]int{124, 58, 237}  // Violet
	colorSuccess   = [3]int{16, 185, 129}  // Green
	colorWarning   = [3]int{245, 158, 11}  // Amber
	colorError     = [3]int{239, 68, 68}   // Red
	colorDark      = [3]int{30, 41, 59}    // Slate 800
	colorMuted     = [3]int{100, 116, 139} // Slate 500
	colorLight     = [3]int{241, 245, 249} // Slate 100
	colorWhite     = [3]int{255, 255, 255}
	colorTableHead = [3]int{51, 65, 85}    // Slate 700
	colorTableAlt  = [3]int{248, 250, 252} // Slate 50
)

type PDFBuilder struct {
	pdf        *gofpdf.Fpdf
	data       *domain.StudentReportData
	pageNumber int
}

func NewPDFBuilder(data *domain.StudentReportData) *PDFBuilder {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(marginLeft, marginTop, marginRight)
	pdf.SetAutoPageBreak(true, marginBottom)

	return &PDFBuilder{pdf: pdf, data: data, pageNumber: 0}
}

func (b *PDFBuilder) Build() ([]byte, error) {
	b.pdf.SetFooterFunc(func() {
		b.pageNumber++
		b.pdf.SetY(-15)
		b.pdf.SetFont("Arial", "", 8)
		b.pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		b.pdf.CellFormat(0, 10, fmt.Sprintf("Pagina %d", b.pageNumber), "", 0, "C", false, 0, "")
	})

	b.pdf.AddPage()
	b.renderHeader()
	b.renderSummaryCards()
	b.renderTopicProgressSection()
	b.renderCourseTable()
	b.renderActivitySection()

	var buf bytes.Buffer
	if err := b.pdf.Output(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (b *PDFBuilder) renderHeader() {
	pdf := b.pdf

	// Top accent bar
	pdf.SetFillColor(colorSecondary[0], colorSecondary[1], colorSecondary[2])
	pdf.Rect(0, 0, pageWidth, 8, "F")

	// Logo/Brand area
	pdf.SetY(15)
	pdf.SetFont("Arial", "B", 24)
	pdf.SetTextColor(colorSecondary[0], colorSecondary[1], colorSecondary[2])
	pdf.CellFormat(contentWidth, 10, "PRACTIQ", "", 1, "L", false, 0, "")

	// Report title
	pdf.SetFont("Arial", "", 11)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.CellFormat(contentWidth, 5, "Reporte de Progreso Academico", "", 1, "L", false, 0, "")
	pdf.Ln(8)

	// Student info card
	pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
	cardY := pdf.GetY()
	pdf.Rect(marginLeft, cardY, contentWidth, 28, "F")

	// Student name
	pdf.SetXY(marginLeft+8, cardY+5)
	pdf.SetFont("Arial", "B", 14)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(100, 6, b.data.Student.Name, "", 0, "L", false, 0, "")

	// Student email
	pdf.SetXY(marginLeft+8, cardY+12)
	pdf.SetFont("Arial", "", 10)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.CellFormat(100, 5, b.data.Student.Email, "", 0, "L", false, 0, "")

	// Period info (right side)
	periodStr := "Periodo: Todo el historial"
	if b.data.Period.From != nil && b.data.Period.To != nil {
		periodStr = fmt.Sprintf("Periodo: %s al %s",
			b.data.Period.From.Format("02/01/2006"),
			b.data.Period.To.Format("02/01/2006"))
	} else if b.data.Period.From != nil {
		periodStr = fmt.Sprintf("Desde: %s", b.data.Period.From.Format("02/01/2006"))
	} else if b.data.Period.To != nil {
		periodStr = fmt.Sprintf("Hasta: %s", b.data.Period.To.Format("02/01/2006"))
	}

	pdf.SetXY(marginLeft+8, cardY+20)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.CellFormat(100, 4, periodStr, "", 0, "L", false, 0, "")

	// Generation date (right side)
	pdf.SetXY(contentWidth-40, cardY+5)
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(40, 5, "Generado:", "", 0, "R", false, 0, "")
	pdf.SetXY(contentWidth-40, cardY+10)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(40, 5, b.data.GeneratedAt.Format("02/01/2006"), "", 0, "R", false, 0, "")
	pdf.SetXY(contentWidth-40, cardY+15)
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
	pdf.CellFormat(40, 5, b.data.GeneratedAt.Format("15:04 hrs"), "", 0, "R", false, 0, "")

	pdf.SetY(cardY + 35)
}

func (b *PDFBuilder) renderSummaryCards() {
	pdf := b.pdf
	s := b.data.Summary

	b.renderSectionTitle("Resumen General")

	// 3 cards per row
	cardWidth := (contentWidth - 8) / 3
	cardHeight := 22.0
	startY := pdf.GetY()

	cards := []struct {
		label string
		value string
		color [3]int
	}{
		{"Temas Practicados", fmt.Sprintf("%d", s.TopicsPracticed), colorPrimary},
		{"Dominio Promedio", fmt.Sprintf("%.1f%%", s.AverageMastery), b.masteryColorRGB(s.AverageMastery)},
		{"Tasa de Acierto", fmt.Sprintf("%.1f%%", s.AccuracyRate), b.masteryColorRGB(s.AccuracyRate)},
		{"Total Intentos", fmt.Sprintf("%d", s.TotalAttempts), colorPrimary},
		{"Respuestas Correctas", fmt.Sprintf("%d", s.CorrectAttempts), colorSuccess},
		{"Racha Actual", fmt.Sprintf("%d dias", s.CurrentStreak), colorSecondary},
	}

	for i, card := range cards {
		col := i % 3
		row := i / 3
		x := marginLeft + float64(col)*(cardWidth+4)
		y := startY + float64(row)*(cardHeight+4)

		// Card background
		pdf.SetFillColor(colorLight[0], colorLight[1], colorLight[2])
		pdf.RoundedRect(x, y, cardWidth, cardHeight, 2, "1234", "F")

		// Color accent bar
		pdf.SetFillColor(card.color[0], card.color[1], card.color[2])
		pdf.Rect(x, y, 3, cardHeight, "F")

		// Value
		pdf.SetXY(x+8, y+4)
		pdf.SetFont("Arial", "B", 16)
		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.CellFormat(cardWidth-12, 8, card.value, "", 0, "L", false, 0, "")

		// Label
		pdf.SetXY(x+8, y+13)
		pdf.SetFont("Arial", "", 8)
		pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		pdf.CellFormat(cardWidth-12, 5, card.label, "", 0, "L", false, 0, "")
	}

	pdf.SetY(startY + 2*(cardHeight+4) + 8)
}

func (b *PDFBuilder) renderTopicProgressSection() {
	if len(b.data.TopicProgress) == 0 {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(80)
	b.renderSectionTitle("Dominio por Tema")

	// Render chart
	chartData, err := GenerateMasteryBarChart(b.data.TopicProgress, 520, 220)
	if err == nil && chartData != nil {
		imgName := "mastery_chart"
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader(imgName, opt, bytes.NewReader(chartData))
		pdf.ImageOptions(imgName, marginLeft, pdf.GetY(), contentWidth, 0, false, opt, 0, "")
		pdf.Ln(58)
	}

	// Top 5 topics table
	b.renderTopTopicsTable()
}

func (b *PDFBuilder) renderTopTopicsTable() {
	pdf := b.pdf
	topics := b.data.TopicProgress

	if len(topics) == 0 {
		return
	}

	// Sort and limit
	limit := 5
	if len(topics) < limit {
		limit = len(topics)
	}

	b.checkPageBreak(50)

	// Table header
	colWidths := []float64{80, 30, 30, 40}
	headers := []string{"Tema", "Dominio", "Nivel", "Intentos"}

	pdf.SetFillColor(colorTableHead[0], colorTableHead[1], colorTableHead[2])
	pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetFont("Arial", "B", 9)

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	for i := 0; i < limit; i++ {
		t := topics[i]

		if i%2 == 0 {
			pdf.SetFillColor(colorTableAlt[0], colorTableAlt[1], colorTableAlt[2])
		} else {
			pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
		}

		// Topic name
		title := t.TopicTitle
		if len(title) > 35 {
			title = title[:32] + "..."
		}
		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.CellFormat(colWidths[0], 7, title, "", 0, "L", true, 0, "")

		// Mastery with color
		mc := b.masteryColorRGB(t.MasteryScore)
		pdf.SetTextColor(mc[0], mc[1], mc[2])
		pdf.CellFormat(colWidths[1], 7, fmt.Sprintf("%.0f%%", t.MasteryScore), "", 0, "C", true, 0, "")

		// Level
		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.CellFormat(colWidths[2], 7, fmt.Sprintf("%d", t.CurrentLevel), "", 0, "C", true, 0, "")

		// Attempts
		pdf.CellFormat(colWidths[3], 7, fmt.Sprintf("%d/%d", t.CorrectAttempts, t.TotalAttempts), "", 0, "C", true, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(8)
}

func (b *PDFBuilder) renderCourseTable() {
	if len(b.data.CourseProgress) == 0 {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(60)
	b.renderSectionTitle("Progreso por Curso")

	// Table header
	colWidths := []float64{70, 25, 30, 30, 25}
	headers := []string{"Curso", "Nivel", "Temas", "Dominio", "Actividad"}

	pdf.SetFillColor(colorTableHead[0], colorTableHead[1], colorTableHead[2])
	pdf.SetTextColor(colorWhite[0], colorWhite[1], colorWhite[2])
	pdf.SetFont("Arial", "B", 9)

	for i, h := range headers {
		pdf.CellFormat(colWidths[i], 8, h, "", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table rows
	pdf.SetFont("Arial", "", 9)
	for i, c := range b.data.CourseProgress {
		if i%2 == 0 {
			pdf.SetFillColor(colorTableAlt[0], colorTableAlt[1], colorTableAlt[2])
		} else {
			pdf.SetFillColor(colorWhite[0], colorWhite[1], colorWhite[2])
		}

		title := c.CourseTitle
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		lastAct := "-"
		if c.LastActivity != nil {
			lastAct = c.LastActivity.Format("02/01")
		}

		pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
		pdf.CellFormat(colWidths[0], 7, title, "", 0, "L", true, 0, "")
		pdf.CellFormat(colWidths[1], 7, fmt.Sprintf("%d", c.CurrentLevel), "", 0, "C", true, 0, "")
		pdf.CellFormat(colWidths[2], 7, fmt.Sprintf("%d", c.TopicCount), "", 0, "C", true, 0, "")

		mc := b.masteryColorRGB(c.AverageMastery)
		pdf.SetTextColor(mc[0], mc[1], mc[2])
		pdf.CellFormat(colWidths[3], 7, fmt.Sprintf("%.0f%%", c.AverageMastery), "", 0, "C", true, 0, "")

		pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])
		pdf.CellFormat(colWidths[4], 7, lastAct, "", 0, "C", true, 0, "")
		pdf.Ln(-1)
	}

	pdf.Ln(8)
}

func (b *PDFBuilder) renderActivitySection() {
	if len(b.data.DailyAttempts) == 0 {
		return
	}

	pdf := b.pdf
	b.checkPageBreak(80)
	b.renderSectionTitle("Actividad Reciente")

	// Render chart
	chartData, err := GenerateDailyAttemptsLineChart(b.data.DailyAttempts, 520, 180)
	if err == nil && chartData != nil {
		imgName := "daily_chart"
		opt := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader(imgName, opt, bytes.NewReader(chartData))
		pdf.ImageOptions(imgName, marginLeft, pdf.GetY(), contentWidth, 0, false, opt, 0, "")
		pdf.Ln(50)
	}

	// Activity stats
	b.renderActivityStats()
}

func (b *PDFBuilder) renderActivityStats() {
	if len(b.data.DailyAttempts) == 0 {
		return
	}

	pdf := b.pdf

	// Calculate stats
	var totalAttempts, totalCorrect, daysActive int
	for _, d := range b.data.DailyAttempts {
		totalAttempts += d.Total
		totalCorrect += d.Correct
		if d.Total > 0 {
			daysActive++
		}
	}

	avgPerDay := 0.0
	if daysActive > 0 {
		avgPerDay = float64(totalAttempts) / float64(daysActive)
	}

	// Stats row
	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(colorMuted[0], colorMuted[1], colorMuted[2])

	statsText := fmt.Sprintf("Dias activos: %d  |  Promedio diario: %.1f intentos  |  Total periodo: %d intentos (%d correctos)",
		daysActive, avgPerDay, totalAttempts, totalCorrect)
	pdf.CellFormat(contentWidth, 6, statsText, "", 1, "C", false, 0, "")
	pdf.Ln(4)
}

func (b *PDFBuilder) renderSectionTitle(title string) {
	pdf := b.pdf

	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(colorDark[0], colorDark[1], colorDark[2])
	pdf.CellFormat(contentWidth, 8, title, "", 1, "L", false, 0, "")

	// Underline
	pdf.SetDrawColor(colorSecondary[0], colorSecondary[1], colorSecondary[2])
	pdf.SetLineWidth(0.5)
	y := pdf.GetY()
	pdf.Line(marginLeft, y, marginLeft+40, y)
	pdf.Ln(6)
}

func (b *PDFBuilder) checkPageBreak(height float64) {
	if b.pdf.GetY()+height > pageHeight-marginBottom-10 {
		b.pdf.AddPage()
	}
}

func (b *PDFBuilder) masteryColorRGB(score float64) [3]int {
	if score >= 80 {
		return colorSuccess
	}
	if score >= 50 {
		return colorWarning
	}
	return colorError
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
}
