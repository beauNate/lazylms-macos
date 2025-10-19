package animation

import (
	"fmt"
	"math"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// grayToForeground generates a color based on position and time for animation
func grayToForeground(i, t float64) lipgloss.Color {
	// Calculate intensity using sine wave
	intensity := (math.Sin(-0.3*i+t) + 1) / 2

	if intensity < 0.5 {
		return styles.ColorWhite
	} else {
		return styles.ColorGray
	}
}

// AnimatedText renders text with animation
func AnimatedText(text string, animationTime time.Time) string {
	var result string
	t := float64(animationTime.UnixMilli()) / 200.0

	for i, ch := range text {
		color := grayToForeground(float64(i), t)
		style := lipgloss.NewStyle().Foreground(color)
		result += style.Render(string(ch))
	}

	return result
}

func GeneratingStatus(animationTime time.Time) string {
	return AnimatedText("generating", animationTime)
}

func AnimatedModelStatus(modelName string, isGenerating bool, animationTime time.Time) string {
	if !isGenerating {
		return modelName
	}

	statusText := fmt.Sprintf("%s - %s", modelName, GeneratingStatus(animationTime))
	return statusText
}
