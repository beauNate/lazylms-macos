package layout

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Rugz007/lazylms/pkg/tui/styles"
)

// Borderize wraps content with a border and optional embedded text at various positions.
// Positions: TopLeft, TopMiddle, TopRight, BottomLeft, BottomMiddle, BottomRight.

type BorderPosition int

const (
	TopLeftBorder BorderPosition = iota
	TopMiddleBorder
	TopRightBorder
	BottomLeftBorder
	BottomMiddleBorder
	BottomRightBorder
)

// Borderize wraps content with a border and optional embedded text at various positions
func Borderize(content string, active bool, width, height int, embeddedText map[BorderPosition]string) string {
	if embeddedText == nil {
		embeddedText = make(map[BorderPosition]string)
	}

	// Constrain content to fit within the border
	content = lipgloss.NewStyle().
		Width(width - 2).
		Height(height - 2).
		Render(content)

	var (
		thickness = map[bool]lipgloss.Border{
			true:  lipgloss.Border(lipgloss.ThickBorder()),
			false: lipgloss.Border(lipgloss.NormalBorder()),
		}
		color = map[bool]lipgloss.TerminalColor{
			true:  styles.ColorPurple,
			false: styles.ColorGray,
		}
		border = thickness[active]
		style  = lipgloss.NewStyle().Foreground(color[active])
	)

	encloseInSquareBrackets := func(text string) string {
		if text != "" {
			return fmt.Sprintf("%s%s%s",
				style.Render(border.TopRight),
				text,
				style.Render(border.TopLeft),
			)
		}
		return text
	}
	buildHorizontalBorder := func(leftText, middleText, rightText, leftCorner, inbetween, rightCorner string) string {
		leftText = encloseInSquareBrackets(leftText)
		middleText = encloseInSquareBrackets(middleText)
		// Don't bracket rightText to avoid corner conflict
		if rightText != "" {
			rightText = style.Render(border.TopRight) + rightText
		}
		// Calculate length of border between embedded texts
		contentWidth := width - 2 // Account for left and right corners
		remaining := max(0, contentWidth-lipgloss.Width(leftText)-lipgloss.Width(middleText)-lipgloss.Width(rightText))
		leftBorderLen := max(0, (contentWidth/2)-lipgloss.Width(leftText)-(lipgloss.Width(middleText)/2))
		rightBorderLen := max(0, remaining-leftBorderLen)
		// Then construct border string
		s := leftText +
			style.Render(strings.Repeat(inbetween, leftBorderLen)) +
			middleText +
			style.Render(strings.Repeat(inbetween, rightBorderLen)) +
			rightText
		// Make it fit in the space available between the two corners.
		s = lipgloss.NewStyle().
			Inline(true).
			MaxWidth(contentWidth).
			Render(s)
		return style.Render(leftCorner) + s + style.Render(rightCorner)
	}
	// Stack top border, content and horizontal borders, and bottom border.
	return strings.Join([]string{
		buildHorizontalBorder(
			embeddedText[TopLeftBorder],
			embeddedText[TopMiddleBorder],
			embeddedText[TopRightBorder],
			border.TopLeft,
			border.Top,
			border.TopRight,
		),
		lipgloss.NewStyle().
			BorderForeground(color[active]).
			Border(border, false, true, false, true).Render(content),
		buildHorizontalBorder(
			embeddedText[BottomLeftBorder],
			embeddedText[BottomMiddleBorder],
			embeddedText[BottomRightBorder],
			border.BottomLeft,
			border.Bottom,
			border.BottomRight,
		),
	}, "\n")
}
