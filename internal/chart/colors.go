package chart

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	DefaultColors   = []string{"#ebedf0", "#9be9a8", "#40c463", "#30a14e", "#216e39"}
	HalloweenColors = []string{"#ebedf0", "#fdf156", "#ffc722", "#ff9500", "#fe6d00"}
	TealColors      = []string{"#ebedf0", "#c6f7d0", "#7fcdcd", "#49a3a3", "#2d7d79"}
)

func GetColorScheme(baseColor string) []string {
	if baseColor == "" {
		return DefaultColors
	}
	
	switch strings.ToLower(baseColor) {
	case "halloween":
		return HalloweenColors
	case "teal":
		return TealColors
	case "default":
		return DefaultColors
	}
	
	if !strings.HasPrefix(baseColor, "#") {
		baseColor = "#" + baseColor
	}
	
	if !isValidHex(baseColor) {
		return DefaultColors
	}
	
	return generateColorScheme(baseColor)
}

func generateColorScheme(baseColor string) []string {
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
	
	if r > 255 { r = 255 }
	if g > 255 { g = 255 }
	if b > 255 { b = 255 }
	
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