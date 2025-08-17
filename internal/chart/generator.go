package chart

import (
	"fmt"
	"strings"
	"time"

	"github.com/xqsit94/ghchart/internal/services"
)

const (
	CellSize = 12
	CellGap  = 3
	Width    = 53 * (CellSize + CellGap) - CellGap
	Height   = 7 * (CellSize + CellGap) - CellGap
	Padding  = 20
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
	data, err := g.githubService.FetchContributions(username)
	if err != nil {
		return nil, err
	}

	colors := GetColorScheme(baseColor)
	svg := g.buildSVG(data, colors)
	
	return []byte(svg), nil
}

func (g *Generator) buildSVG(data *services.ContributionYear, colors []string) string {
	totalWidth := Width + 2*Padding
	totalHeight := Height + 2*Padding

	var svg strings.Builder
	
	svg.WriteString(fmt.Sprintf(`<svg width="%d" height="%d" xmlns="http://www.w3.org/2000/svg">`, totalWidth, totalHeight))
	svg.WriteString(fmt.Sprintf(`<rect width="%d" height="%d" fill="white"/>`, totalWidth, totalHeight))
	svg.WriteString(fmt.Sprintf(`<g transform="translate(%d,%d)">`, Padding, Padding))

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

	currentDate := startDate
	week := 0

	for week < 53 {
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
				`<rect x="%d" y="%d" width="%d" height="%d" fill="%s" rx="2" ry="2">`,
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