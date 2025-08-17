package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ContributionData struct {
	Date         string
	Count        int
	Level        int
	Color        string
}

type ContributionYear struct {
	Year          int
	Contributions []ContributionData
	Total         int
}

func fetchGitHubContributions(username string) (*ContributionYear, error) {
	// Use GitHub's contributions endpoint
	url := fmt.Sprintf("https://github.com/users/%s/contributions", username)
	
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	// Add headers to mimic a real browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	
	resp, err := client.Do(req)
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

	// Extract contribution data from td elements with data-date attributes
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
		
		// Try to extract count from tooltip or other attributes
		// GitHub sometimes puts count info in tool-tip elements or aria-label
		count := 0
		if toolTip := s.Find("tool-tip"); toolTip.Length() > 0 {
			// Extract count from tooltip content if available
			if ariaLabel, exists := toolTip.Attr("aria-label"); exists {
				count = extractCountFromTooltip(ariaLabel)
			}
		}
		
		// If we don't have a count but have a level, estimate based on level
		if count == 0 && level > 0 {
			// GitHub uses levels 0-4, so we can estimate:
			// Level 1: 1-3 contributions, Level 2: 4-6, Level 3: 7-9, Level 4: 10+
			switch level {
			case 1:
				count = 2
			case 2:
				count = 5
			case 3:
				count = 8
			case 4:
				count = 12
			}
		}

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

// extractCountFromTooltip tries to extract contribution count from tooltip text
func extractCountFromTooltip(tooltip string) int {
	// Look for patterns like "5 contributions on Jan 1, 2024"
	if strings.Contains(tooltip, "contribution") {
		parts := strings.Fields(tooltip)
		if len(parts) > 0 {
			if count, err := strconv.Atoi(parts[0]); err == nil {
				return count
			}
		}
	}
	return 0
}

func removeExtension(filename, ext string) string {
	if strings.HasSuffix(filename, ext) {
		return strings.TrimSuffix(filename, ext)
	}
	return filename
}