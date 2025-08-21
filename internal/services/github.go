package services

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ContributionData struct {
	Date  string
	Count int
	Level int
}

type ContributionYear struct {
	Year          int
	Contributions []ContributionData
	Total         int
}

type GitHubService struct {
	client *http.Client
}

func NewGitHubService() *GitHubService {
	return &GitHubService{
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (g *GitHubService) FetchContributions(username string) (*ContributionYear, error) {
	url := fmt.Sprintf("https://github.com/users/%s/contributions", username)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to fetch GitHub contributions: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	year := &ContributionYear{
		Year:          time.Now().Year(),
		Contributions: []ContributionData{},
		Total:         0,
	}

	doc.Find("td[data-date]").Each(func(i int, s *goquery.Selection) {
		date, exists := s.Attr("data-date")
		if !exists {
			return
		}

		levelStr, exists := s.Attr("data-level")
		if !exists {
			return
		}

		level, _ := strconv.Atoi(levelStr)
		count := g.estimateCountFromLevel(level)

		contribution := ContributionData{
			Date:  date,
			Count: count,
			Level: level,
		}

		year.Contributions = append(year.Contributions, contribution)
		year.Total += count
	})

	if len(year.Contributions) == 0 {
		return nil, fmt.Errorf("no contribution data found for user %s", username)
	}

	return year, nil
}

func (g *GitHubService) estimateCountFromLevel(level int) int {
	switch level {
	case 1:
		return 2
	case 2:
		return 5
	case 3:
		return 8
	case 4:
		return 12
	default:
		return 0
	}
}

func RemoveExtension(filename, ext string) string {
	if strings.HasSuffix(filename, ext) {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}
