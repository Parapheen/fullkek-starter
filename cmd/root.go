package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Parapheen/fullkek-starter/internal/tui/output"
)

// RootCommand constructs the base command for the fullkek CLI.
func RootCommand() *cobra.Command {
	var bannerPrinted bool

	cmd := &cobra.Command{
		Use:   "fullkek",
		Short: "Scaffold hypermedia-driven Go applications with opinionated stacks.",
		Long: `Fullkek bootstraps hypermedia-first web applications written in Go.

Choose an opinionated stack—like Go + HTMX + Tailwind CSS or Go + DataStar + Bootstrap—and
build a project skeleton that is ready for templates, assets, and server wiring.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			if bannerPrinted {
				return
			}
			fmt.Fprintln(cmd.OutOrStdout(), output.RenderBanner())
			fmt.Fprintln(cmd.OutOrStdout())
			bannerPrinted = true
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().BoolP("verbose", "v", false, "enable verbose output")

	cmd.AddCommand(newNewCommand())

	return cmd
}
