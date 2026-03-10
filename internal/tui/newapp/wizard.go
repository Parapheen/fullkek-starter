package newapp

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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

type multiSelectBinding struct {
	category stacks.FeatureCategory
	choices  []stacks.Feature
	values   []string
}

const sqliteFeatureID = "database-sqlite"

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

	workingDir, err := os.Getwd()
	if err != nil {
		workingDir = ""
	}

	if outputDir == "" {
		outputDir = suggestOutputDir(appName)
	}

	totalSteps := len(opts.Categories) + 5
	step := 1

	groups := make([]*huh.Group, 0, 3+len(opts.Categories)+2)

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title(stepLabel(step, totalSteps, "Application name")).
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
	step++

	groups = append(groups, huh.NewGroup(
		huh.NewInput().
			Title(stepLabel(step, totalSteps, "Go module path")).
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
	step++

	bindings := make([]*featureBinding, 0, len(opts.Categories))
	multiBindings := make([]*multiSelectBinding, 0)

	var databaseBinding *featureBinding
	var authBinding *featureBinding
	authEnabled := func() bool {
		if databaseBinding == nil {
			return false
		}
		return strings.TrimSpace(databaseBinding.value) == sqliteFeatureID
	}
	oauth2Selected := func() bool {
		if authBinding == nil {
			return false
		}
		return strings.TrimSpace(authBinding.value) == "auth-oauth2"
	}

	for _, category := range opts.Categories {
		choices := opts.FeatureChoices[category.ID]

		if category.AllowMultiple && category.ID == stacks.CategoryOAuthProviders {
			// Multi-select for OAuth providers
			mBinding := &multiSelectBinding{category: category, choices: choices}
			defaults := defaultSelection[category.ID]
			if len(defaults) > 0 {
				mBinding.values = defaults
			}

			options := make([]huh.Option[string], 0, len(choices))
			for _, feature := range choices {
				options = append(options, huh.NewOption(feature.Name, feature.ID))
			}

			multiSelect := huh.NewMultiSelect[string]().
				Title(stepLabel(step, totalSteps, category.Name)).
				Options(options...).
				Value(&mBinding.values)

			if description := strings.TrimSpace(category.Description); description != "" {
				multiSelect.Description(description)
			}

			multiBindings = append(multiBindings, mBinding)

			group := huh.NewGroup(multiSelect)
			group.WithHideFunc(func() bool {
				return !oauth2Selected()
			})

			groups = append(groups, group)
			step++
			continue
		}

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
			Title(stepLabel(step, totalSteps, category.Name)).
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
		if category.ID == stacks.CategoryDatabase {
			databaseBinding = binding
		}
		if category.ID == stacks.CategoryAuth {
			authBinding = binding
		}

		group := huh.NewGroup(selectField)
		if category.ID == stacks.CategoryAuth {
			group.WithHideFunc(func() bool {
				return !authEnabled()
			})
		}
		// Hide email/payments/deploy categories that have "none" as default
		// unless DB is SQLite (for payments)
		if category.ID == stacks.CategoryEmail {
			group.WithHideFunc(func() bool {
				return !authEnabled()
			})
		}
		if category.ID == stacks.CategoryPayments {
			group.WithHideFunc(func() bool {
				return !authEnabled()
			})
		}

		groups = append(groups, group)
		step++
	}

	groups = append(groups, huh.NewGroup(
		huh.NewConfirm().
			Title(stepLabel(step, totalSteps, "Overwrite destination if it already exists?")).
			DescriptionFunc(func() string {
				destination := scaffoldDestinationPath(workingDir, appName, outputDir)
				if destination == "" {
					return "Press enter to toggle."
				}
				return fmt.Sprintf("The project will be scaffolded at %s.\nPress enter to toggle.", destination)
			}, struct {
				AppName   *string
				OutputDir *string
			}{
				AppName:   &appName,
				OutputDir: &outputDir,
			}).
			Affirmative("Yes").
			Negative("No").
			Value(&force),
	))
	step++

	scaffoldNow := true

	buildSummary := func() string {
		trimmedApp := strings.TrimSpace(appName)
		trimmedModule := strings.TrimSpace(modulePath)
		resolvedOutput := resolveOutputDir(appName, outputDir)
		if strings.TrimSpace(resolvedOutput) == "" {
			resolvedOutput = "(not set)"
		}
		destination := scaffoldDestinationPath(workingDir, appName, outputDir)
		if destination == "" {
			destination = "(not set)"
		}

		var lines []string
		lines = append(lines, "Project details")
		lines = append(lines,
			fmt.Sprintf("  App name    : %s", valueOrPlaceholder(trimmedApp)),
			fmt.Sprintf("  Module path : %s", valueOrPlaceholder(trimmedModule)),
			fmt.Sprintf("  Output dir  : %s", resolvedOutput),
			fmt.Sprintf("  Destination : %s", destination),
			fmt.Sprintf("  Overwrite   : %s", humanizeBool(force)),
		)

		selected := make(map[string]string, len(bindings))
		for _, binding := range bindings {
			if binding.category.ID == stacks.CategoryAuth && !authEnabled() {
				continue
			}
			selected[binding.category.ID] = binding.value
		}

		selection := stacks.SelectionFromIDs(selected)
		// Add multi-select values
		for _, mb := range multiBindings {
			if oauth2Selected() && len(mb.values) > 0 {
				selection[mb.category.ID] = mb.values
			}
		}

		stack, err := stacks.Compose(selection)
		if err == nil {
			lines = append(lines, "")
			lines = append(lines, "Stack")
			lines = append(lines, fmt.Sprintf("  %s", valueOrPlaceholder(strings.TrimSpace(stack.Name))))
		} else {
			lines = append(lines, "")
			lines = append(lines, "Stack")
			lines = append(lines, fmt.Sprintf("  Invalid selection: %s", err.Error()))
		}

		featureBlocksAdded := false
		for _, binding := range bindings {
			if binding.category.ID == stacks.CategoryAuth && !authEnabled() {
				continue
			}
			if strings.TrimSpace(binding.value) == "" {
				continue
			}

			categoryName := strings.TrimSpace(binding.category.Name)
			if categoryName == "" {
				categoryName = "Other"
			}

			if !featureBlocksAdded {
				lines = append(lines, "")
				lines = append(lines, "Selections")
				featureBlocksAdded = true
			}

			lines = append(lines, fmt.Sprintf("  %s: %s", categoryName, binding.value))
		}
		for _, mb := range multiBindings {
			if oauth2Selected() && len(mb.values) > 0 {
				if !featureBlocksAdded {
					lines = append(lines, "")
					lines = append(lines, "Selections")
					featureBlocksAdded = true
				}
				lines = append(lines, fmt.Sprintf("  %s: %s", mb.category.Name, strings.Join(mb.values, ", ")))
			}
		}

		lines = append(lines, "")
		lines = append(lines, "After generation")
		lines = append(lines, fmt.Sprintf("  cd %s", resolvedOutput))
		lines = append(lines, "  make go")

		return strings.Join(lines, "\n")
	}

	reviewNote := huh.NewNote()
	reviewNote = reviewNote.
		Title(stepLabel(step, totalSteps, "Review configuration")).
		DescriptionFunc(func() string {
			return buildSummary()
		}, struct {
			AppName    *string
			ModulePath *string
			OutputDir  *string
			Force      *bool
		}{
			AppName:    &appName,
			ModulePath: &modulePath,
			OutputDir:  &outputDir,
			Force:      &force,
		}).
		Height(16)

	groups = append(groups, huh.NewGroup(reviewNote))
	step++

	groups = append(groups, huh.NewGroup(
		huh.NewConfirm().
			Title(stepLabel(step, totalSteps, "Ready to scaffold?")).
			Description("Use Shift+Tab to edit previous answers, or choose Cancel.").
			Affirmative("Scaffold now").
			Negative("Cancel").
			Value(&scaffoldNow),
	))

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

	if !scaffoldNow {
		return Result{}, ErrCancelled
	}

	appName = strings.TrimSpace(appName)
	modulePath = strings.TrimSpace(modulePath)
	outputDir = strings.TrimSpace(outputDir)
	if outputDir == "" {
		outputDir = suggestOutputDir(appName)
	}

	authAvailable := authEnabled()

	selection := make(map[string]string, len(bindings))
	for _, binding := range bindings {
		if binding.category.ID == stacks.CategoryAuth && !authAvailable {
			continue
		}
		selection[binding.category.ID] = binding.value
	}

	result := stacks.SelectionFromIDs(selection)

	// Add multi-select OAuth providers
	for _, mb := range multiBindings {
		if oauth2Selected() && len(mb.values) > 0 {
			result[mb.category.ID] = mb.values
		}
	}

	return Result{
		AppName:    appName,
		ModulePath: modulePath,
		OutputDir:  outputDir,
		Selection:  result,
		Force:      force,
	}, nil
}

func suggestOutputDir(appName string) string {
	return resolveOutputDir(appName, "")
}

func first(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func scaffoldDestinationPath(baseDir, appName, outputDir string) string {
	destination := resolveOutputDir(appName, outputDir)
	destination = strings.TrimSpace(destination)
	if destination == "" {
		return ""
	}
	if baseDir != "" && !filepath.IsAbs(destination) {
		destination = filepath.Join(baseDir, destination)
	}
	return filepath.Clean(destination)
}

func resolveOutputDir(appName, override string) string {
	override = strings.TrimSpace(override)
	if override != "" {
		return override
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

func valueOrPlaceholder(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "(not set)"
	}
	return value
}

func humanizeBool(value bool) string {
	if value {
		return "Yes"
	}
	return "No"
}

func stepLabel(step, total int, title string) string {
	return fmt.Sprintf("Step %d/%d: %s", step, total, title)
}
