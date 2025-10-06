package scaffold

import "github.com/Parapheen/fullkek-starter/internal/stacks"

// BaseDirectories are created for every generated project regardless of stack.
var BaseDirectories = []string{
	"cmd/server",
	"internal/app",
	"internal/transport/http",
	"web/templates/pages",
	"web/assets/styles",
	"web/assets/scripts",
	"bin",
}

// BaseTemplates are rendered for every generated project.
var BaseTemplates = []stacks.Template{
	{
		Source:      "base/go.mod.tmpl",
		Destination: "go.mod",
		Mode:        0o644,
	},
	{
		Source:      "base/README.md.tmpl",
		Destination: "README.md",
		Mode:        0o644,
	},
	{
		Source:      "base/Makefile.tmpl",
		Destination: "Makefile",
		Mode:        0o644,
	},
	{
		Source:      "base/.gitignore.tmpl",
		Destination: ".gitignore",
		Mode:        0o644,
	},
	{
		Source:      "base/cmd/server/main.go.tmpl",
		Destination: "cmd/server/main.go",
		Mode:        0o644,
	},
	{
		Source:      "base/internal/app/app.go.tmpl",
		Destination: "internal/app/app.go",
		Mode:        0o644,
	},
	{
		Source:      "base/internal/transport/http/server.go.tmpl",
		Destination: "internal/transport/http/server.go",
		Mode:        0o644,
	},
	{
		Source:      "base/internal/transport/http/router.go.tmpl",
		Destination: "internal/transport/http/router.go",
		Mode:        0o644,
	},
	{
		Source:      "base/web/templates/README.md.tmpl",
		Destination: "web/templates/README.md",
		Mode:        0o644,
	},
	{
		Source:      "base/web/assets/styles/README.md.tmpl",
		Destination: "web/assets/styles/README.md",
		Mode:        0o644,
	},
	{
		Source:      "base/web/assets/scripts/README.md.tmpl",
		Destination: "web/assets/scripts/README.md",
		Mode:        0o644,
	},
}
