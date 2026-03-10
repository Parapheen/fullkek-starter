package cmd

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Parapheen/fullkek-starter/internal/scaffold"
	"github.com/Parapheen/fullkek-starter/internal/stacks"
	"github.com/Parapheen/fullkek-starter/internal/tui/newapp"
	"github.com/Parapheen/fullkek-starter/internal/tui/output"
)

func newNewCommand() *cobra.Command {
	defaults := stacks.DefaultSelection()

	var opts struct {
		modulePath     string
		outputDir      string
		force          bool
		noUI           bool
		frontend       string
		styling        string
		http           string
		database       string
		auth           string
		oauthProviders string
		email          string
		payments       string
		deploy         string
	}

	frontendDefault := first(defaults[stacks.CategoryFrontend])
	stylingDefault := first(defaults[stacks.CategoryStyling])
	httpDefault := first(defaults[stacks.CategoryHTTP])
	databaseDefault := first(defaults[stacks.CategoryDatabase])
	authDefault := first(defaults[stacks.CategoryAuth])
	emailDefault := first(defaults[stacks.CategoryEmail])
	paymentsDefault := first(defaults[stacks.CategoryPayments])
	deployDefault := first(defaults[stacks.CategoryDeploy])

	cmd := &cobra.Command{
		Use:   "new [app-name]",
		Short: "Create a new hypermedia project using modular feature blocks.",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var appName string
			if len(args) > 0 {
				appName = args[0]
			}

			flagSelection := stacks.MergeSelections(
				stacks.DefaultSelection(),
				stacks.SelectionFromIDs(map[string]string{
					stacks.CategoryFrontend: opts.frontend,
					stacks.CategoryStyling:  opts.styling,
					stacks.CategoryHTTP:     opts.http,
					stacks.CategoryDatabase: opts.database,
					stacks.CategoryAuth:     opts.auth,
					stacks.CategoryEmail:    opts.email,
					stacks.CategoryPayments: opts.payments,
					stacks.CategoryDeploy:   opts.deploy,
				}),
			)

			// Parse comma-separated OAuth providers
			if opts.oauthProviders != "" {
				providers := strings.Split(opts.oauthProviders, ",")
				providerIDs := make([]string, 0, len(providers))
				for _, p := range providers {
					p = strings.TrimSpace(p)
					if p != "" {
						if !strings.HasPrefix(p, "oauth-") {
							p = "oauth-" + p
						}
						providerIDs = append(providerIDs, p)
					}
				}
				if len(providerIDs) > 0 {
					flagSelection[stacks.CategoryOAuthProviders] = providerIDs
				}
			}

			selection := stacks.CloneSelection(flagSelection)

			var (
				modulePath  string
				destination string
				force       = opts.force
			)

			categories := stacks.Categories()

			if shouldUseWizard(cmd, opts.noUI) {
				featureChoices := make(map[string][]stacks.Feature, len(categories))
				for _, category := range categories {
					featureChoices[category.ID] = stacks.FeaturesForCategory(category.ID)
				}

				wizardResult, err := newapp.Run(newapp.Options{
					AppName:          appName,
					ModulePath:       opts.modulePath,
					OutputDir:        opts.outputDir,
					Force:            opts.force,
					Categories:       categories,
					FeatureChoices:   featureChoices,
					DefaultSelection: stacks.CloneSelection(selection),
				}, cmd.InOrStdin(), cmd.OutOrStdout())
				if err != nil {
					if errors.Is(err, newapp.ErrCancelled) {
						return nil
					}
					return err
				}

				selection = stacks.MergeSelections(selection, wizardResult.Selection)
				appName = strings.TrimSpace(wizardResult.AppName)
				modulePath = strings.TrimSpace(wizardResult.ModulePath)
				destination = strings.TrimSpace(wizardResult.OutputDir)
				force = wizardResult.Force
			} else {
				if appName == "" {
					return errors.New("app name required when not using interactive mode")
				}
				modulePath = deriveModulePath(appName, opts.modulePath)
				destination = opts.outputDir
			}

			appName = strings.TrimSpace(appName)
			modulePath = strings.TrimSpace(modulePath)
			destination = strings.TrimSpace(destination)

			if appName == "" {
				return errors.New("app name cannot be empty")
			}
			if modulePath == "" {
				modulePath = deriveModulePath(appName, modulePath)
			}
			destination = deriveOutputDir(appName, destination)

			stack, err := stacks.Compose(selection)
			if err != nil {
				return err
			}

			generator := scaffold.DefaultGenerator()
			ctx := context.Background()

			if verbose(cmd) {
				fmt.Fprintf(cmd.ErrOrStderr(), "Scaffolding %s at %s using %s\n", appName, destination, stack.Name)
			}

			if err := generator.Generate(ctx, scaffold.Options{
				AppName:     appName,
				ModulePath:  modulePath,
				Destination: destination,
				Stack:       stack,
				Force:       force,
			}); err != nil {
				return err
			}

			if shouldUseWizard(cmd, opts.noUI) {
				output.PrintSuccess(cmd.OutOrStdout(), destination, stack)
			} else {
				printNextSteps(cmd.OutOrStdout(), destination, stack)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&opts.modulePath, "module", "", "Go module path for the generated project")
	cmd.Flags().StringVarP(&opts.outputDir, "output", "o", "", "target directory (defaults to app name)")
	cmd.Flags().BoolVar(&opts.force, "force", false, "overwrite destination directory if it already exists")
	cmd.Flags().BoolVar(&opts.noUI, "no-ui", false, "disable the interactive wizard")
	cmd.Flags().StringVar(&opts.frontend, "frontend", frontendDefault, "frontend runtime feature identifier")
	cmd.Flags().StringVar(&opts.styling, "styling", stylingDefault, "styling feature identifier")
	cmd.Flags().StringVar(&opts.http, "http", httpDefault, "HTTP framework feature identifier")
	cmd.Flags().StringVar(&opts.database, "database", databaseDefault, "database feature identifier")
	cmd.Flags().StringVar(&opts.auth, "auth", authDefault, "authentication feature identifier")
	cmd.Flags().StringVar(&opts.oauthProviders, "oauth-providers", "", "comma-separated OAuth providers (github,google,yandex)")
	cmd.Flags().StringVar(&opts.email, "email", emailDefault, "email sending feature identifier")
	cmd.Flags().StringVar(&opts.payments, "payments", paymentsDefault, "payment processing feature identifier")
	cmd.Flags().StringVar(&opts.deploy, "deploy", deployDefault, "deployment feature identifier")

	registerFeatureCompletion(cmd, "frontend", stacks.CategoryFrontend)
	registerFeatureCompletion(cmd, "styling", stacks.CategoryStyling)
	registerFeatureCompletion(cmd, "http", stacks.CategoryHTTP)
	registerFeatureCompletion(cmd, "database", stacks.CategoryDatabase)
	registerFeatureCompletion(cmd, "auth", stacks.CategoryAuth)
	registerFeatureCompletion(cmd, "email", stacks.CategoryEmail)
	registerFeatureCompletion(cmd, "payments", stacks.CategoryPayments)
	registerFeatureCompletion(cmd, "deploy", stacks.CategoryDeploy)

	return cmd
}

func deriveModulePath(appName, override string) string {
	if override != "" {
		return override
	}
	sanitized := strings.TrimSpace(appName)
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ToLower(sanitized)
	return sanitized
}

func deriveOutputDir(appName, override string) string {
	if override != "" {
		return strings.TrimSpace(override)
	}
	sanitized := strings.TrimSpace(appName)
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ToLower(sanitized)
	sanitized = strings.Trim(sanitized, "-")
	if sanitized == "" {
		return strings.TrimSpace(appName)
	}
	return sanitized
}

func isInteractive(r io.Reader) bool {
	file, ok := r.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func isInteractiveWriter(w io.Writer) bool {
	file, ok := w.(*os.File)
	if !ok {
		return false
	}
	info, err := file.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}

func shouldUseWizard(cmd *cobra.Command, skip bool) bool {
	if skip {
		return false
	}
	return isInteractive(cmd.InOrStdin()) && isInteractiveWriter(cmd.OutOrStdout())
}

func verbose(cmd *cobra.Command) bool {
	v, _ := cmd.Flags().GetBool("verbose")
	return v
}

func printNextSteps(out io.Writer, destination string, stack stacks.Stack) {
	fmt.Fprintf(out, "\nProject scaffolded at %s\n", destination)
	fmt.Fprintf(out, "Stack: %s\n", stack.Name)

	if len(stack.Tags) > 0 {
		fmt.Fprintf(out, "Tags: %s\n", strings.Join(stack.Tags, ", "))
	}
	if len(stack.Features) > 0 {
		fmt.Fprintln(out, "Features:")
		for _, feature := range stack.Features {
			fmt.Fprintf(out, "  - %s\n", feature.Name)
		}
	}

	fmt.Fprintln(out, "\nNext steps:")
	fmt.Fprintf(out, "  1. cd %s\n", destination)
	fmt.Fprintln(out, "  2. Review .env and fill any required credentials")
	fmt.Fprintln(out, "  3. make go")
	fmt.Fprintln(out, "  4. Open http://localhost:3333")
	fmt.Fprintf(out, "\nReview %s/README.md for detailed guidance.\n", destination)
}

func registerFeatureCompletion(cmd *cobra.Command, flagName, categoryID string) {
	err := cmd.RegisterFlagCompletionFunc(flagName, func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		features := stacks.FeaturesForCategory(categoryID)
		values := make([]string, 0, len(features))
		for _, feature := range features {
			values = append(values, feature.ID+"\t"+feature.Name)
		}
		return values, cobra.ShellCompDirectiveNoFileComp
	})
	if err != nil {
		panic(err)
	}
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
