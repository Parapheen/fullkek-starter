package assets

import (
	_ "embed"
	"strings"
)

// HorseBanner stores the ASCII art horse banner for reuse across the CLI.
//
//go:embed horse.txt
var HorseBanner string

// HorseBannerLines returns the horse banner split into individual lines.
func HorseBannerLines() []string {
	trimmed := strings.TrimRight(HorseBanner, "\n")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "\n")
}
