package main

import (
	"fmt"
	"strings"
	"time"
)

func generateChart(username, baseColor string) ([]byte, error) {
	data, err := fetchGitHubContributions(username)
	if err != nil {
		return nil, err
	}

	colors := getColorScheme(baseColor)
	svg := buildSVG(data, colors)
	
	return []byte(svg), nil
}

func buildSVG(data *ContributionYear, colors []string) string {
	const (
		cellSize = 12
		cellGap  = 3
		width    = 53 * (cellSize + cellGap) - cellGap  // 53 weeks
		height   = 7 * (cellSize + cellGap) - cellGap   // 7 days
		padding  = 20
	)

	totalWidth := width + 2*padding
	totalHeight := height + 2*padding

	var svg strings.Builder
	
	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalHeight))
	svg.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, totalWidth, totalHeight))
	svg.WriteString(fmt.Sprintf(`<g transform="translate(%d,%d)">`, padding, padding))

	// Create a map of dates to contribution data
	dateMap := make(map[string]ContributionData)
	for _, contrib := range data.Contributions {
		dateMap[contrib.Date] = contrib
	}

	// Get the last day of contributions (most recent)
	endDate := time.Now()
	if len(data.Contributions) > 0 {
		if lastDate, err := time.Parse("2006-01-02", data.Contributions[len(data.Contributions)-1].Date); err == nil {
			endDate = lastDate
		}
	}

	// Start from 53 weeks ago (approximately 1 year)
	startDate := endDate.AddDate(0, 0, -364) // 52 weeks * 7 days = 364 days

	// Align to start of week (Sunday)
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	currentDate := startDate
	week := 0

	for week < 53 {
		for day := 0; day < 7; day++ {
			dateStr := currentDate.Format("2006-01-02")
			
			x := week * (cellSize + cellGap)
			y := day * (cellSize + cellGap)
			
			color := colors[0] // Default no contribution color
			
			if contrib, exists := dateMap[dateStr]; exists {
				if contrib.Level >= 0 && contrib.Level < len(colors) {
					color = colors[contrib.Level]
				}
			}
			
			svg.WriteString(fmt.Sprintf(
				`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="2" ry="2">`,
				x, y, cellSize, cellSize, color,
			))
			svg.WriteString(fmt.Sprintf(
				`<title>%s: %d contributions</title>`,
				dateStr, getContributionCount(dateMap, dateStr),
			))
			svg.WriteString(`</rect>`)
			
			currentDate = currentDate.AddDate(0, 0, 1)
		}
		week++
	}

	svg.WriteString(`</g>`)
	svg.WriteString(`</svg>`)

	return svg.String()
}

func getContributionCount(dateMap map[string]ContributionData, date string) int {
	if contrib, exists := dateMap[date]; exists {
		return contrib.Count
	}
	return 0
}