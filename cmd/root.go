package cmd

import (
	"github.com/spf13/cobra"
)

// RootCommand constructs the base command for the fullkek CLI.
func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fullkek",
		Short: "Scaffold hypermedia-driven Go applications with opinionated stacks.",
		Long: `Fullkek bootstraps hypermedia-first web applications written in Go.

Choose an opinionated stack—like Go + HTMX + Tailwind CSS or Go + DataStar + Bootstrap—and
build a project skeleton that is ready for templates, assets, and server wiring.`,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	cmd.AddCommand(newNewCommand())

	return cmd
}
