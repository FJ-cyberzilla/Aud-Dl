package ui

import (
	"fmt"
	"strings"
)

type BarChart struct {
	Title string
	Data  map[string]float64
	Theme Theme
	Width int
}

func NewBarChart(title string, data map[string]float64, theme Theme) *BarChart {
	return &BarChart{
		Title: title,
		Data:  data,
		Theme: theme,
		Width: 30,
	}
}

func (bc *BarChart) Render() string {
	var sb strings.Builder
	sb.WriteString(TrueColor(fmt.Sprintf("\n📊 %s\n", bc.Title), bc.Theme.Primary))

	var maxVal float64
	for _, val := range bc.Data {
		if val > maxVal {
			maxVal = val
		}
	}

	for label, val := range bc.Data {
		barLen := int((val / maxVal) * float64(bc.Width))
		bar := strings.Repeat("█", barLen) + strings.Repeat("░", bc.Width-barLen)

		line := fmt.Sprintf("  %-15s │ %s %.1f%%\n", label, TrueColor(bar, bc.Theme.Secondary), val)
		sb.WriteString(line)
	}

	return sb.String()
}
