package newapp

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/huh"

	"github.com/Parapheen/fullkek-starter/internal/stacks"
)

// ErrCancelled signals that the user aborted the wizard.
var ErrCancelled = errors.New("wizard cancelled")

// Options seeds the wizard with defaults.
type Options struct {
	AppName          string
	ModulePath       string
	OutputDir        string
	Force            bool
	Categories       []stacks.FeatureCategory
	FeatureChoices   map[string][]stacks.Feature
	DefaultSelection stacks.Selection
}

// Result captures the selections produced by the wizard.
type Result struct {
	AppName    string
	ModulePath string
	OutputDir  string
	Selection  stacks.Selection
	Force      bool
}

type featureBinding struct {
	category stacks.FeatureCategory
	choices  []stacks.Feature
	value    string
}

// Run executes the wizard and returns the user's selections.
func Run(opts Options, input io.Reader, output io.Writer) (Result, error) {
	if len(opts.Categories) == 0 {
		return Result{}, errors.New("no feature categories available for interactive wizard")
	}

	for _, category := range opts.Categories {
		if len(opts.FeatureChoices[category.ID]) == 0 {
			return Result{}, fmt.Errorf("no features registered for category %q", category.Name)
		}
	}

	defaultSelection := stacks.CloneSelection(opts.DefaultSelection)

	appName := strings.TrimSpace(opts.AppName)
	modulePath := strings.TrimSpace(opts.ModulePath)
	outputDir := strings.TrimSpace(opts.OutputDir)
	force := opts.Force

	if outputDir == "" {
		outputDir = suggestOutputDir(appName)
	}

	groups := make([]*huh.Group, 0, 3+len(opts.Categories)+2)

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Application name").
			Description("Used to derive default module and destination names.").
			Placeholder("fullkek-demo").
			Value(&appName).
			Validate(func(value string) error {
				if strings.TrimSpace(value) == "" {
					return errors.New("app name cannot be empty")
				}
				return nil
			}),
	))

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title("Go module path").
			Description("Enter the Go module path (e.g. github.com/username/project).").
			Placeholder("github.com/username/project").
			Value(&modulePath).
			Validate(func(value string) error {
				if strings.TrimSpace(value) == "" {
					return errors.New("module path cannot be empty")
				}
				return nil
			}),
	))

	bindings := make([]*featureBinding, 0, len(opts.Categories))

	for _, category := range opts.Categories {
		choices := opts.FeatureChoices[category.ID]
		binding := &featureBinding{category: category, choices: choices}

		defaultID := first(defaultSelection[category.ID])
		if defaultID != "" {
			for _, feature := range choices {
				if feature.ID == defaultID {
					binding.value = defaultID
					break
				}
			}
		}
		if binding.value == "" {
			binding.value = choices[0].ID
		}

		options := make([]huh.Option[string], 0, len(choices))
		for _, feature := range choices {
			options = append(options, huh.NewOption(feature.Name, feature.ID))
		}

		selectField := huh.NewSelect[string]().
			Title(category.Name).
			Options(options...).
			Value(&binding.value).
			Validate(func(id string) error {
				if strings.TrimSpace(id) == "" {
					return fmt.Errorf("select a feature for %s", category.Name)
				}
				return nil
			})

		if description := strings.TrimSpace(category.Description); description != "" {
			selectField.Description(description)
		}

		bindings = append(bindings, binding)
		groups = append(groups, huh.NewGroup(selectField))
	}

	groups = append(groups, huh.NewGroup(
		huh.NewConfirm().
			Title("Overwrite destination if it already exists?").
			Description("Press enter to toggle.").
			Affirmative("Yes").
			Negative("No").
			Value(&force),
	))

	summary := huh.NewNote().
		Title("Review configuration").
		DescriptionFunc(func() string {
			var b strings.Builder
			writeLine := func(format string, args ...any) {
				fmt.Fprintf(&b, format+"\n", args...)
			}

			trimmedApp := strings.TrimSpace(appName)
			trimmedModule := strings.TrimSpace(modulePath)
			trimmedOutput := strings.TrimSpace(outputDir)

			writeLine("App name      : %s", trimmedApp)
			writeLine("Module path   : %s", trimmedModule)
			writeLine("Destination   : %s", trimmedOutput)
			if len(bindings) > 0 {
				b.WriteString("\nFeatures:\n")
				for _, binding := range bindings {
					featureName := binding.selectedFeatureName()
					if featureName != "" && featureName != "<none>" {
						writeLine("  • %s: %s", binding.category.Name, featureName)
					}
				}
			}
			writeLine("\nForce overwrite: %t", force)
			b.WriteString("\n✨ Press Enter to scaffold or use ← to adjust previous answers.")
			return b.String()
		}, nil).
		Next(true).
		NextLabel("🚀 Scaffold project")

	groups = append(groups, huh.NewGroup(summary))

	form := huh.NewForm(groups...)

	if input != nil {
		form = form.WithInput(input)
	}
	if output != nil {
		form = form.WithOutput(output)
	}

	if err := form.Run(); err != nil {
		if errors.Is(err, huh.ErrUserAborted) {
			return Result{}, ErrCancelled
		}
		return Result{}, err
	}

	appName = strings.TrimSpace(appName)
	modulePath = strings.TrimSpace(modulePath)
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		outputDir = suggestOutputDir(appName)
	}

	selection := make(map[string]string, len(bindings))
	for _, binding := range bindings {
		selection[binding.category.ID] = binding.value
	}

	return Result{
		AppName:    appName,
		ModulePath: modulePath,
		OutputDir:  outputDir,
		Selection:  stacks.SelectionFromIDs(selection),
		Force:      force,
	}, nil
}

func (b *featureBinding) selectedFeatureName() string {
	for _, feature := range b.choices {
		if feature.ID == b.value {
			return feature.Name
		}
	}
	return "<none>"
}

func suggestOutputDir(appName string) string {
	sanitized := strings.TrimSpace(appName)
	if sanitized == "" {
		return ""
	}
	sanitized = strings.ReplaceAll(sanitized, " ", "-")
	sanitized = strings.ToLower(sanitized)
	sanitized = strings.Trim(sanitized, "-")
	return sanitized
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}
