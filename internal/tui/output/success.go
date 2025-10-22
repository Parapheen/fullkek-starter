package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/Parapheen/fullkek-starter/internal/stacks"
)

var (
	// Color palette
	primaryColor    = lipgloss.Color("#00D9FF")
	successColor    = lipgloss.Color("#00FF87")
	accentColor     = lipgloss.Color("#FF00FF")
	mutedColor      = lipgloss.Color("#808080")
	warningColor    = lipgloss.Color("#FFD700")
	backgroundColor = lipgloss.Color("#1a1a1a")

	// Styles
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(successColor).
			MarginTop(1).
			MarginBottom(1)

	boxStyle = lipgloss.NewStyle().
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1)

	infoStyle = lipgloss.NewStyle().
			Foreground(primaryColor)

	labelStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Bold(true)

	valueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	stepStyle = lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true)

	featureStyle = lipgloss.NewStyle().
			Foreground(accentColor)

	bannerLines = []string{
		"███████╗██╗   ██╗██╗     ██╗     ██╗  ██╗███████╗██╗  ██╗",
		"██╔════╝██║   ██║██║     ██║     ██║ ██╔╝██╔════╝██║ ██╔╝",
		"█████╗  ██║   ██║██║     ██║     █████╔╝ █████╗  █████╔╝ ",
		"██╔══╝  ██║   ██║██║     ██║     ██╔═██╗ ██╔══╝  ██╔═██╗ ",
		"██║     ╚██████╔╝███████╗███████╗██║  ██╗███████╗██║  ██╗",
		"╚═╝      ╚═════╝ ╚══════╝╚══════╝╚═╝  ╚═╝╚══════╝╚═╝  ╚═╝",
	}

	bannerGradient = []lipgloss.Color{
		lipgloss.Color("#FF5F6D"),
		lipgloss.Color("#FF7E79"),
		lipgloss.Color("#FF9E73"),
		lipgloss.Color("#FFC66F"),
		lipgloss.Color("#F7E96E"),
		lipgloss.Color("#D9FF7E"),
		lipgloss.Color("#7CE8FF"),
		lipgloss.Color("#5CC2FF"),
		lipgloss.Color("#4A95FF"),
		lipgloss.Color("#826CFF"),
		lipgloss.Color("#B15CFF"),
		lipgloss.Color("#DF5CFF"),
	}

	bannerBaseStyle = lipgloss.NewStyle().
			Bold(true)

	bannerContainerStyle = lipgloss.NewStyle().
				Background(backgroundColor).
				Padding(2, 2)
)

// RenderBanner returns the styled fullkek header with its gradient applied.
func RenderBanner() string {
	styledLines := make([]string, len(bannerLines))
	for i, line := range bannerLines {
		color := bannerGradient[i%len(bannerGradient)]
		styledLines[i] = bannerBaseStyle.Foreground(color).Render(line)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, styledLines...)
	return bannerContainerStyle.Render(content)
}

// PrintSuccess renders a beautiful success message with project details and next steps
func PrintSuccess(out io.Writer, destination string, stack stacks.Stack) {
	var b strings.Builder

	// Next steps
	var steps strings.Builder
	steps.WriteString(infoStyle.Render("🚀 Next steps:") + "\n\n")
	steps.WriteString(fmt.Sprintf("   %s %s\n", stepStyle.Render("1."), valueStyle.Render(fmt.Sprintf("cd %s", destination))))
	steps.WriteString(fmt.Sprintf("   %s %s\n", stepStyle.Render("2."), valueStyle.Render("make go")))

	mutedStyle := lipgloss.NewStyle().Foreground(mutedColor)
	steps.WriteString(fmt.Sprintf("\n   %s %s\n", labelStyle.Render("📖"), mutedStyle.Render(fmt.Sprintf("Review %s/README.md for detailed guidance.", destination))))

	stepsBox := boxStyle.Render(steps.String())
	b.WriteString(stepsBox + "\n")

	fmt.Fprint(out, b.String())
}
