package styles

import (
	"github.com/charmbracelet/lipgloss"
)

// Color constants
const (
	ColorWhite     = lipgloss.Color("15")      // White
	ColorPurple    = lipgloss.Color("63")      // Purple
	ColorPurpleBg  = lipgloss.Color("62")      // Purple background
	ColorYellow    = lipgloss.Color("#f0c024") // Yellow/light
	ColorGray      = lipgloss.Color("241")     // Gray
	ColorBlack     = lipgloss.Color("0")       // Black
	ColorGreen     = lipgloss.Color("2")       // Green
	ColorBlue      = lipgloss.Color("4")       // Blue
	ColorOrange    = lipgloss.Color("208")     // Orange
	ColorReasoning = lipgloss.Color("245")     // Dimmed gray for reasoning text
)

// Adaptive colors (can't be const)
var (
	ColorForeground         = lipgloss.AdaptiveColor{Light: string(ColorBlack), Dark: string(ColorWhite)} // Black in light mode, white in dark mode
	ColorForegroundInverted = lipgloss.AdaptiveColor{Light: string(ColorWhite), Dark: string(ColorBlack)}
)

// Theme defines a color theme for the application
type Theme struct {
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Info      lipgloss.Color
}

// DefaultTheme returns the default color theme
func DefaultTheme() Theme {
	return Theme{
		Primary:   ColorPurple,
		Secondary: ColorGray,
		Success:   ColorGreen,
		Warning:   ColorYellow,
		Error:     ColorOrange,
		Info:      ColorBlue,
	}
}
