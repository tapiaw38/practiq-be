package studentreport

import (
	"bytes"
	"sort"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

var (
	colorPrimary   = drawing.Color{R: 124, G: 58, B: 237, A: 255}  // violet
	colorSecondary = drawing.Color{R: 99, G: 102, B: 241, A: 255}  // indigo
	colorSuccess   = drawing.Color{R: 34, G: 197, B: 94, A: 255}   // green
	colorWarning   = drawing.Color{R: 234, G: 179, B: 8, A: 255}   // yellow
	colorError     = drawing.Color{R: 239, G: 68, B: 68, A: 255}   // red
	colorGray      = drawing.Color{R: 156, G: 163, B: 175, A: 255} // gray
)

func masteryColor(score float64) drawing.Color {
	if score >= 80 {
		return colorSuccess
	}
	if score >= 50 {
		return colorWarning
	}
	return colorError
}

func GenerateMasteryBarChart(topics []domain.StudentTopicProgress, width, height int) ([]byte, error) {
	if len(topics) == 0 {
		return nil, nil
	}

	// Sort by mastery descending and take top 10
	sorted := make([]domain.StudentTopicProgress, len(topics))
	copy(sorted, topics)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].MasteryScore > sorted[j].MasteryScore
	})
	if len(sorted) > 10 {
		sorted = sorted[:10]
	}

	var bars []chart.Value
	for _, t := range sorted {
		label := t.TopicTitle
		if len(label) > 20 {
			label = label[:17] + "..."
		}
		bars = append(bars, chart.Value{
			Label: label,
			Value: t.MasteryScore,
			Style: chart.Style{
				FillColor:   masteryColor(t.MasteryScore),
				StrokeColor: masteryColor(t.MasteryScore),
				StrokeWidth: 0,
			},
		})
	}

	barChart := chart.BarChart{
		Title:      "Dominio por Tema (Top 10)",
		TitleStyle: chart.Style{FontSize: 12, FontColor: drawing.ColorBlack},
		Width:      width,
		Height:     height,
		BarWidth:   30,
		XAxis: chart.Style{
			FontSize: 8,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{Min: 0, Max: 100},
			Style: chart.Style{FontSize: 9},
		},
		Bars: bars,
	}

	buf := bytes.NewBuffer(nil)
	if err := barChart.Render(chart.PNG, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func GenerateDailyAttemptsLineChart(daily []domain.DailyAttemptCount, width, height int) ([]byte, error) {
	if len(daily) == 0 {
		return nil, nil
	}

	// Sort by date
	sorted := make([]domain.DailyAttemptCount, len(daily))
	copy(sorted, daily)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Date.Before(sorted[j].Date)
	})

	var xValues []float64
	var totalValues []float64
	var correctValues []float64

	for i, d := range sorted {
		xValues = append(xValues, float64(i))
		totalValues = append(totalValues, float64(d.Total))
		correctValues = append(correctValues, float64(d.Correct))
	}

	lineChart := chart.Chart{
		Title:      "Intentos por Dia",
		TitleStyle: chart.Style{FontSize: 12, FontColor: drawing.ColorBlack},
		Width:      width,
		Height:     height,
		XAxis: chart.XAxis{
			Style: chart.Style{FontSize: 8},
			ValueFormatter: func(v interface{}) string {
				idx := int(v.(float64))
				if idx >= 0 && idx < len(sorted) {
					return sorted[idx].Date.Format("02/01")
				}
				return ""
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{FontSize: 9},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Total",
				XValues: xValues,
				YValues: totalValues,
				Style: chart.Style{
					StrokeColor: colorPrimary,
					StrokeWidth: 2,
				},
			},
			chart.ContinuousSeries{
				Name:    "Correctas",
				XValues: xValues,
				YValues: correctValues,
				Style: chart.Style{
					StrokeColor: colorSuccess,
					StrokeWidth: 2,
				},
			},
		},
	}

	lineChart.Elements = []chart.Renderable{
		chart.Legend(&lineChart),
	}

	buf := bytes.NewBuffer(nil)
	if err := lineChart.Render(chart.PNG, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
