package chart

import (
	"fmt"
	"strings"
	"time"

	"github.com/xqsit94/ghchart/internal/services"
)

const (
	CellSize    = 12
	CellGap     = 3
	Width       = 53*(CellSize+CellGap) - CellGap
	Height      = 7*(CellSize+CellGap) - CellGap
	Padding     = 15
	LabelOffset = 15
)

type Generator struct {
	githubService *services.GitHubService
}

func NewGenerator(githubService *services.GitHubService) *Generator {
	return &Generator{
		githubService: githubService,
	}
}

func (g *Generator) Generate(username, baseColor string) ([]byte, error) {
	return g.GenerateWithTheme(username, baseColor, "light")
}

func (g *Generator) GenerateWithTheme(username, baseColor, theme string) ([]byte, error) {
	data, err := g.githubService.FetchContributions(username)
	if err != nil {
		return nil, err
	}

	colors := GetColorSchemeWithTheme(baseColor, theme)
	svg := g.buildSVG(data, colors, theme)

	return []byte(svg), nil
}

func (g *Generator) buildSVG(data *services.ContributionYear, colors []string, theme string) string {
	totalWidth := Width + Padding + LabelOffset
	totalHeight := Height + LabelOffset

	var svg strings.Builder

	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalHeight))

	labelColor := "#1f2328"
	if strings.ToLower(theme) == "dark" {
		labelColor = "#cecece"
	}

	svg.WriteString(`<style>`)
	svg.WriteString(fmt.Sprintf(`.month-label { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 12px; fill: %s; }`, labelColor))
	svg.WriteString(fmt.Sprintf(`.day-label { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; font-size: 12px; fill: %s; }`, labelColor))
	svg.WriteString(`</style>`)

	svg.WriteString(fmt.Sprintf(`<g transform="translate(%d,%d)">`, Padding+LabelOffset, LabelOffset))

	dateMap := make(map[string]services.ContributionData)
	for _, contrib := range data.Contributions {
		dateMap[contrib.Date] = contrib
	}

	endDate := time.Now()
	if len(data.Contributions) > 0 {
		if lastDate, err := time.Parse("2006-01-02", data.Contributions[len(data.Contributions)-1].Date); err == nil {
			endDate = lastDate
		}
	}

	startDate := endDate.AddDate(0, 0, -364)
	for startDate.Weekday() != time.Sunday {
		startDate = startDate.AddDate(0, 0, -1)
	}

	dayLabels := []string{"", "Mon", "", "Wed", "", "Fri", ""}
	for i, label := range dayLabels {
		if label != "" {
			y := i*(CellSize+CellGap) + CellSize/2
			svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" text-anchor="end" class="day-label" dominant-baseline="middle">%s</text>`,
				-5, y, label))
		}
	}

	currentDate := startDate
	week := 0
	var monthPositions []struct {
		name string
		x    int
	}

	prevMonth := ""

	for week < 53 {
		weekStartDate := currentDate
		currentMonth := weekStartDate.Format("Jan")
		if currentMonth != prevMonth && week > 0 {
			x := week * (CellSize + CellGap)
			monthPositions = append(monthPositions, struct {
				name string
				x    int
			}{currentMonth, x})
			prevMonth = currentMonth
		} else if week == 0 {
			x := week * (CellSize + CellGap)
			monthPositions = append(monthPositions, struct {
				name string
				x    int
			}{currentMonth, x})
			prevMonth = currentMonth
		}

		for day := 0; day < 7; day++ {
			dateStr := currentDate.Format("2006-01-02")

			x := week * (CellSize + CellGap)
			y := day * (CellSize + CellGap)

			color := colors[0]

			if contrib, exists := dateMap[dateStr]; exists {
				if contrib.Level >= 0 && contrib.Level < len(colors) {
					color = colors[contrib.Level]
				}
			}

			svg.WriteString(fmt.Sprintf(
				`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" stroke="rgba(27,31,36,0.06)" stroke-width="1" rx="2" ry="2">`,
				x, y, CellSize, CellSize, color,
			))
			svg.WriteString(fmt.Sprintf(
				`<title>%s: %d contributions</title>`,
				dateStr, g.getContributionCount(dateMap, dateStr),
			))
			svg.WriteString(`</rect>`)

			currentDate = currentDate.AddDate(0, 0, 1)
		}
		week++
	}

	for _, pos := range monthPositions {
		svg.WriteString(fmt.Sprintf(`<text x="%d" y="%d" class="month-label">%s</text>`,
			pos.x, -5, pos.name))
	}

	svg.WriteString(`</g>`)
	svg.WriteString(`</svg>`)

	return svg.String()
}

func (g *Generator) getContributionCount(dateMap map[string]services.ContributionData, date string) int {
	if contrib, exists := dateMap[date]; exists {
		return contrib.Count
	}
	return 0
}
