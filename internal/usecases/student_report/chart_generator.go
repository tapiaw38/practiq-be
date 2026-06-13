package studentreport

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/tapiaw38/practiq-be/internal/domain"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
)

var (
	chartColorPrimary   = drawing.Color{R: 99, G: 102, B: 241, A: 255}  // Indigo
	chartColorSecondary = drawing.Color{R: 124, G: 58, B: 237, A: 255}  // Violet
	chartColorSuccess   = drawing.Color{R: 16, G: 185, B: 129, A: 255}  // Green
	chartColorWarning   = drawing.Color{R: 245, G: 158, B: 11, A: 255}  // Amber
	chartColorError     = drawing.Color{R: 239, G: 68, B: 68, A: 255}   // Red
	chartColorMuted     = drawing.Color{R: 148, G: 163, B: 184, A: 255} // Slate 400
	chartColorDark      = drawing.Color{R: 51, G: 65, B: 85, A: 255}    // Slate 700
	chartColorLight     = drawing.Color{R: 241, G: 245, B: 249, A: 255} // Slate 100
	chartColorBg        = drawing.Color{R: 255, G: 255, B: 255, A: 255}
)

func masteryColor(score float64) drawing.Color {
	if score >= 80 {
		return chartColorSuccess
	}
	if score >= 50 {
		return chartColorWarning
	}
	return chartColorError
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
		if len(label) > 18 {
			label = label[:15] + "..."
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
		Title: "",
		Background: chart.Style{
			FillColor: chartColorBg,
			Padding: chart.Box{
				Top:    20,
				Left:   10,
				Right:  10,
				Bottom: 40,
			},
		},
		Canvas: chart.Style{
			FillColor: chartColorBg,
		},
		Width:    width,
		Height:   height,
		BarWidth: 35,
		BarSpacing: 8,
		XAxis: chart.Style{
			FontSize:  8,
			FontColor: chartColorDark,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{Min: 0, Max: 100},
			Style: chart.Style{
				FontSize:  9,
				FontColor: chartColorMuted,
			},
			ValueFormatter: func(v interface{}) string {
				return fmt.Sprintf("%.0f%%", v.(float64))
			},
			GridMajorStyle: chart.Style{
				StrokeColor: chartColorLight,
				StrokeWidth: 1,
			},
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

	// Limit to last 14 days for cleaner chart
	if len(sorted) > 14 {
		sorted = sorted[len(sorted)-14:]
	}

	var xValues []float64
	var totalValues []float64
	var correctValues []float64

	for i, d := range sorted {
		xValues = append(xValues, float64(i))
		totalValues = append(totalValues, float64(d.Total))
		correctValues = append(correctValues, float64(d.Correct))
	}

	lineChart := chart.Chart{
		Title: "",
		Background: chart.Style{
			FillColor: chartColorBg,
			Padding: chart.Box{
				Top:    15,
				Left:   10,
				Right:  15,
				Bottom: 30,
			},
		},
		Canvas: chart.Style{
			FillColor: chartColorBg,
		},
		Width:  width,
		Height: height,
		XAxis: chart.XAxis{
			Style: chart.Style{
				FontSize:  8,
				FontColor: chartColorMuted,
			},
			ValueFormatter: func(v interface{}) string {
				idx := int(v.(float64))
				if idx >= 0 && idx < len(sorted) {
					return sorted[idx].Date.Format("02/01")
				}
				return ""
			},
			GridMajorStyle: chart.Style{
				StrokeColor: chartColorLight,
				StrokeWidth: 1,
			},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{
				FontSize:  9,
				FontColor: chartColorMuted,
			},
			GridMajorStyle: chart.Style{
				StrokeColor: chartColorLight,
				StrokeWidth: 1,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Total",
				XValues: xValues,
				YValues: totalValues,
				Style: chart.Style{
					StrokeColor: chartColorPrimary,
					StrokeWidth: 2.5,
					DotColor:    chartColorPrimary,
					DotWidth:    4,
				},
			},
			chart.ContinuousSeries{
				Name:    "Correctas",
				XValues: xValues,
				YValues: correctValues,
				Style: chart.Style{
					StrokeColor: chartColorSuccess,
					StrokeWidth: 2.5,
					DotColor:    chartColorSuccess,
					DotWidth:    4,
				},
			},
		},
	}

	lineChart.Elements = []chart.Renderable{
		chart.LegendThin(&lineChart),
	}

	buf := bytes.NewBuffer(nil)
	if err := lineChart.Render(chart.PNG, buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
