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
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
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

	memeStyle = lipgloss.NewStyle().
			Foreground(primaryColor).
			MarginTop(1).
			MarginBottom(1)
)

const coolMeme = `
  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣿⡿⠿⠛⠛⠛⠛⠛⠛⠿⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⡿⠁⠀⠀  👓  ⠀⠀⠀⠈⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⡿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⡇⠀⠀  FULLKEK  ⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⡇⠀⠀  DEPLOY   ⠀⠀⠀⢸⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣧⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣼⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣷⣄⠀⠀⠀⠀⠀⠀⠀⠀⣠⣾⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣶⣤⣀⣀⣤⣶⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿
  ⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿

          IT'S ALIVE! ⚡️ Time to ship! 🚀
`

// PrintSuccess renders a beautiful success message with project details and next steps
func PrintSuccess(out io.Writer, destination string, stack stacks.Stack) {
	var b strings.Builder

	// Meme
	// memeText := memeStyle.Render(coolMeme)
	// b.WriteString(memeText + "\n")

	// Project details box
	var details strings.Builder
	details.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("📦 Project:"), valueStyle.Render(destination)))
	details.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("🏗️  Stack:"), valueStyle.Render(stack.Name)))

	if len(stack.Tags) > 0 {
		details.WriteString(fmt.Sprintf("%s %s\n", labelStyle.Render("🏷️  Tags:"), valueStyle.Render(strings.Join(stack.Tags, ", "))))
	}

	if len(stack.Features) > 0 {
		details.WriteString(fmt.Sprintf("\n%s\n", labelStyle.Render("✨ Features:")))
		for _, feature := range stack.Features {
			details.WriteString(fmt.Sprintf("   %s %s\n", featureStyle.Render("•"), valueStyle.Render(feature.Name)))
		}
	}

	boxContent := boxStyle.Render(details.String())
	b.WriteString(boxContent + "\n")

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
