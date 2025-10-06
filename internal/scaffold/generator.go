package scaffold

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/Parapheen/fullkek-starter/internal/stacks"
	"github.com/Parapheen/fullkek-starter/internal/templates"
)

// Options describes the parameters required to scaffold a project.
type Options struct {
	AppName     string
	ModulePath  string
	Destination string
	Stack       stacks.Stack
	Force       bool
}

// Generator renders the templates embedded in the CLI.
type Generator struct {
	fs fs.FS
}

// NewGenerator returns a generator instance backed by the embedded templates.
func NewGenerator(fsys fs.FS) *Generator {
	return &Generator{fs: fsys}
}

// Generate scaffolds the project according to the provided options.
func (g *Generator) Generate(ctx context.Context, opts Options) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if opts.AppName == "" {
		return errors.New("app name is required")
	}
	if opts.ModulePath == "" {
		return errors.New("module path is required")
	}

	root := opts.Destination
	if root == "" {
		root = opts.AppName
	}

	if err := ensureDestination(root, opts.Force); err != nil {
		return err
	}

	for _, dir := range append(BaseDirectories, opts.Stack.Directories...) {
		if err := ctx.Err(); err != nil {
			return err
		}
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0o755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	data := struct {
		AppName    string
		ModulePath string
		Stack      stacks.Stack
		Generated  time.Time
	}{
		AppName:    opts.AppName,
		ModulePath: opts.ModulePath,
		Stack:      opts.Stack,
		Generated:  time.Now().UTC(),
	}

	templatesToRender := append([]stacks.Template{}, BaseTemplates...)
	templatesToRender = append(templatesToRender, opts.Stack.Templates...)

	for _, tmpl := range templatesToRender {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := g.renderTemplate(root, tmpl, data); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) renderTemplate(root string, tmpl stacks.Template, data any) error {
	target := filepath.Join(root, tmpl.Destination)
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return fmt.Errorf("prepare directory for %s: %w", tmpl.Destination, err)
	}

	funcMap := template.FuncMap{
		"has": func(needle string, haystack []string) bool {
			for _, item := range haystack {
				if item == needle {
					return true
				}
			}
			return false
		},
	}

	parsed, err := template.New(filepath.Base(tmpl.Source)).Funcs(funcMap).Option("missingkey=error").ParseFS(g.fs, tmpl.Source)
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmpl.Source, err)
	}

	file, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, fileModeOrDefault(tmpl.Mode))
	if err != nil {
		return fmt.Errorf("open target %s: %w", tmpl.Destination, err)
	}
	defer file.Close()

	if err := parsed.Execute(file, data); err != nil {
		return fmt.Errorf("execute template %s: %w", tmpl.Source, err)
	}

	return nil
}

func ensureDestination(path string, force bool) error {
	info, err := os.Stat(path)
	if err == nil {
		if !info.IsDir() {
			return fmt.Errorf("destination %s exists and is not a directory", path)
		}
		if !force {
			entries, readErr := os.ReadDir(path)
			if readErr != nil {
				return fmt.Errorf("inspect destination %s: %w", path, readErr)
			}
			if len(entries) > 0 {
				return fmt.Errorf("destination %s is not empty (use --force to overwrite)", path)
			}
		}
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("checking destination %s: %w", path, err)
	}

	return os.MkdirAll(path, 0o755)
}

func fileModeOrDefault(mode fs.FileMode) fs.FileMode {
	if mode == 0 {
		return 0o644
	}
	return mode
}

// DefaultGenerator instantiates a generator with the embedded templates.
func DefaultGenerator() *Generator {
	return NewGenerator(templates.Files)
}
