package chart

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	DefaultColors       = []string{"#ebedf0", "#9be9a8", "#40c463", "#30a14e", "#216e39"}
	DefaultColorsDark   = []string{"#0d1116", "#0e4429", "#006d32", "#26a641", "#39d353"}
	HalloweenColors     = []string{"#ebedf0", "#fdf156", "#ffc722", "#ff9500", "#fe6d00"}
	HalloweenColorsDark = []string{"#0d1116", "#ffdf5d", "#ffa500", "#ff6b35", "#e63946"}
	TealColors          = []string{"#ebedf0", "#c6f7d0", "#7fcdcd", "#49a3a3", "#2d7d79"}
	TealColorsDark      = []string{"#0d1116", "#17a2b8", "#20c997", "#28a745", "#6f42c1"}
)

type ThemeConfig struct {
	Background string
	Border     string
	Text       string
}

var (
	LightTheme = ThemeConfig{
		Background: "#ffffff",
		Border:     "#d1d9e0",
		Text:       "#24292f",
	}
	DarkTheme = ThemeConfig{
		Background: "#0d1117",
		Border:     "#30363d",
		Text:       "#f0f6fc",
	}
)

func GetThemeConfig(theme string) ThemeConfig {
	switch strings.ToLower(theme) {
	case "dark":
		return DarkTheme
	case "light", "":
		return LightTheme
	default:
		return LightTheme
	}
}

func ParseThemeColor(themeColorParam string) (theme, color string) {
	if strings.Contains(themeColorParam, ":") {
		parts := strings.SplitN(themeColorParam, ":", 2)
		return parts[0], parts[1]
	}
	return "light", themeColorParam
}

func GetColorScheme(baseColor string) []string {
	return GetColorSchemeWithTheme(baseColor, "light")
}

func GetColorSchemeWithTheme(baseColor, theme string) []string {
	isDark := strings.ToLower(theme) == "dark"

	if baseColor == "" {
		if isDark {
			return DefaultColorsDark
		}
		return DefaultColors
	}

	switch strings.ToLower(baseColor) {
	case "halloween":
		if isDark {
			return HalloweenColorsDark
		}
		return HalloweenColors
	case "teal":
		if isDark {
			return TealColorsDark
		}
		return TealColors
	case "default":
		if isDark {
			return DefaultColorsDark
		}
		return DefaultColors
	}

	if !strings.HasPrefix(baseColor, "#") {
		baseColor = "#" + baseColor
	}

	if !isValidHex(baseColor) {
		if isDark {
			return DefaultColorsDark
		}
		return DefaultColors
	}

	return generateColorScheme(baseColor, isDark)
}

func generateColorScheme(baseColor string, isDark bool) []string {
	if isDark {
		return []string{
			"#0d1116",
			darkenColor(baseColor, 0.6),
			darkenColor(baseColor, 0.3),
			baseColor,
			lightenColor(baseColor, 0.3),
		}
	}
	return []string{
		"#ebedf0",
		lightenColor(baseColor, 0.7),
		baseColor,
		darkenColor(baseColor, 0.8),
		darkenColor(baseColor, 0.6),
	}
}

func lightenColor(hexColor string, amount float64) string {
	r, g, b := hexToRGB(hexColor)

	r = int(float64(r) + (255-float64(r))*amount)
	g = int(float64(g) + (255-float64(g))*amount)
	b = int(float64(b) + (255-float64(b))*amount)

	if r > 255 {
		r = 255
	}
	if g > 255 {
		g = 255
	}
	if b > 255 {
		b = 255
	}

	return rgbToHex(r, g, b)
}

func darkenColor(hexColor string, amount float64) string {
	r, g, b := hexToRGB(hexColor)

	r = int(float64(r) * amount)
	g = int(float64(g) * amount)
	b = int(float64(b) * amount)

	return rgbToHex(r, g, b)
}

func hexToRGB(hex string) (int, int, int) {
	hex = strings.TrimPrefix(hex, "#")

	if len(hex) != 6 {
		return 0, 0, 0
	}

	r, _ := strconv.ParseInt(hex[0:2], 16, 0)
	g, _ := strconv.ParseInt(hex[2:4], 16, 0)
	b, _ := strconv.ParseInt(hex[4:6], 16, 0)

	return int(r), int(g), int(b)
}

func rgbToHex(r, g, b int) string {
	return fmt.Sprintf("#%02x%02x%02x", r, g, b)
}

func isValidHex(hex string) bool {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return false
	}

	_, err := strconv.ParseInt(hex, 16, 0)
	return err == nil
}
