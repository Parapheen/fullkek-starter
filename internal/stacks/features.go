package stacks

import "sort"

const (
	CategoryFrontend = "frontend"
	CategoryStyling  = "styling"
	CategoryHTTP     = "http"
	CategoryDatabase = "database"
	CategoryAuth     = "auth"
)

var categories = []FeatureCategory{
	{
		ID:          CategoryFrontend,
		Name:        "Frontend runtime",
		Description: "Choose the hypermedia enhancement layer.",
		Required:    true,
	},
	{
		ID:          CategoryStyling,
		Name:        "Styling system",
		Description: "Pick the preferred CSS framework or utility approach.",
		Required:    true,
	},
	{
		ID:          CategoryHTTP,
		Name:        "Web framework",
		Description: "Choose the HTTP framework powering the transport.",
		Required:    true,
	},
	{
		ID:            CategoryDatabase,
		Name:          "Database",
		Description:   "Select the database adapter for persistence needs.",
		AllowMultiple: false,
	},
	{
		ID:            CategoryAuth,
		Name:          "Authentication",
		Description:   "Optional authentication providers.",
		AllowMultiple: false,
	},
}

var featureCatalog = []Feature{
	{
		ID:          "frontend-htmx",
		CategoryID:  CategoryFrontend,
		Name:        "HTMX",
		Description: "Server-driven interactions with HTMX requests and swaps.",
		Tags:        []string{"HTMX"},
		Directories: []string{
			"public/assets/scripts",
		},
		Templates: []Template{
			{
				Source:      "features/frontend/htmx/web/templates/pages/home.html.tmpl",
				Destination: "web/templates/pages/home.html",
			},
			{
				Source:      "features/frontend/htmx/assets/scripts/htmx.min.js.tmpl",
				Destination: "public/assets/scripts/htmx.min.js",
			},
		},
	},
	{
		ID:          "frontend-fixi",
		CategoryID:  CategoryFrontend,
		Name:        "Fixi.js",
		Description: "Composable DOM bindings using the Fixi.js micro-library.",
		Tags:        []string{"Fixi.js"},
		Templates: []Template{
			{
				Source:      "features/frontend/fixi/web/templates/pages/home.html.tmpl",
				Destination: "web/templates/pages/home.html",
			},
			{
				Source:      "features/frontend/fixi/public/assets/scripts/fixi.js.tmpl",
				Destination: "public/assets/scripts/fixi.js",
			},
		},
	},
	{
		ID:          "styling-tailwind",
		CategoryID:  CategoryStyling,
		Name:        "Tailwind CSS",
		Description: "Utility-first styling powered by standalone Tailwind CLI binary.",
		Tags:        []string{"Tailwind"},
		Directories: []string{
			"web/assets/styles/tokens",
		},
		Templates: []Template{
			{
				Source:      "features/styling/tailwind/web/assets/styles/input.css.tmpl",
				Destination: "web/assets/styles/input.css",
			},
			{
				Source:      "features/styling/tailwind/public/assets/styles/output.css.tmpl",
				Destination: "public/assets/styles/output.css",
			},
			{
				Source:      "features/styling/tailwind/web/templates/pages/index.html.tmpl",
				Destination: "web/templates/pages/index.html",
			},
		},
	},
	{
		ID:          "styling-tailwind-basecoat",
		CategoryID:  CategoryStyling,
		Name:        "Tailwind CSS + Basecoat",
		Description: "Tailwind standalone CLI with Basecoat component library via CDN.",
		Tags:        []string{"Tailwind", "Basecoat"},
		Directories: []string{
			"web/assets/styles/tokens",
		},
		Templates: []Template{
			{
				Source:      "features/styling/tailwind_basecoat/web/assets/styles/input.css.tmpl",
				Destination: "web/assets/styles/input.css",
			},
			{
				Source:      "features/styling/tailwind_basecoat/public/assets/styles/output.css.tmpl",
				Destination: "public/assets/styles/output.css",
			},
			{
				Source:      "features/styling/tailwind_basecoat/web/templates/pages/index.html.tmpl",
				Destination: "web/templates/pages/index.html",
			},
		},
	},
	{
		ID:          "styling-daisyui",
		CategoryID:  CategoryStyling,
		Name:        "DaisyUI standalone",
		Description: "Tailwind standalone CLI plus DaisyUI fast script generated bundle.",
		Tags:        []string{"DaisyUI"},
		Directories: []string{
			"public/assets/styles",
		},
		Templates: []Template{
			{
				Source:      "features/styling/daisyui/public/assets/styles/custom.css.tmpl",
				Destination: "public/assets/styles/custom.css",
			},
			{
				Source:      "features/styling/daisyui/public/assets/styles/output.css.tmpl",
				Destination: "public/assets/styles/output.css",
			},
			{
				Source:      "features/styling/daisyui/web/assets/styles/input.css.tmpl",
				Destination: "web/assets/styles/input.css",
			},
			{
				Source:      "features/styling/daisyui/web/templates/pages/index.html.tmpl",
				Destination: "web/templates/pages/index.html",
			},
		},
	},
	{
		ID:          "http-standard",
		CategoryID:  CategoryHTTP,
		Name:        "net/http",
		Description: "Standard library HTTP server with a ServeMux and HTML response.",
		Tags:        []string{"net/http"},
		Templates: []Template{
			{
				Source:      "features/http/standard/internal/transport/http/server.go.tmpl",
				Destination: "internal/transport/http/server.go",
			},
			{
				Source:      "features/http/standard/internal/transport/http/router.go.tmpl",
				Destination: "internal/transport/http/router.go",
			},
		},
	},
	{
		ID:          "http-chi",
		CategoryID:  CategoryHTTP,
		Name:        "Chi",
		Description: "Go-chi router with middleware-ready structure.",
		Tags:        []string{"chi"},
		Templates: []Template{
			{
				Source:      "features/http/chi/internal/transport/http/server.go.tmpl",
				Destination: "internal/transport/http/server.go",
			},
			{
				Source:      "features/http/chi/internal/transport/http/router.go.tmpl",
				Destination: "internal/transport/http/router.go",
			},
		},
	},
	{
		ID:          "database-none",
		CategoryID:  CategoryDatabase,
		Name:        "None",
		Description: "Skip bundling a database integration.",
		Tags:        []string{"database"},
	},
	{
		ID:          "database-sqlite",
		CategoryID:  CategoryDatabase,
		Name:        "SQLite",
		Description: "Preconfigured SQLite helper powered by sqlx.",
		Tags:        []string{"database", "SQLite", "sqlx"},
		Directories: []string{
			"internal/infrastructure/persistence",
		},
		Templates: []Template{
			{
				Source:      "features/database/sqlite/internal/infrastructure/persistence/sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/sqlite.go",
			},
		},
	},
	{
		ID:          "auth-github-oauth2",
		CategoryID:  CategoryAuth,
		Name:        "GitHub (oauth2)",
		Description: "Login with GitHub using x/oauth2 and server-side sessions.",
		Tags:        []string{"auth", "oauth2", "github"},
		Directories: []string{
			"db/migrations",
			"internal/app/auth",
			"internal/domain/session",
			"internal/domain/user",
			"internal/infrastructure/auth",
			"internal/infrastructure/http",
			"internal/infrastructure/persistence",
			"internal/transport/http",
			"web/templates/pages",
		},
		Templates: []Template{
			// HTTP transport
			{
				Source:      "features/auth/github-oauth2/internal/transport/http/oauth_handlers.go.tmpl",
				Destination: "internal/transport/http/oauth_handlers.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/transport/http/render.go.tmpl",
				Destination: "internal/transport/http/render.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/transport/http/auth_middleware.go.tmpl",
				Destination: "internal/transport/http/auth_middleware.go",
			},
			// Web templates
			{
				Source:      "features/auth/github-oauth2/web/templates/pages/profile.html.tmpl",
				Destination: "web/templates/pages/profile.html",
			},
			{
				Source:      "features/auth/github-oauth2/web/templates/pages/login.html.tmpl",
				Destination: "web/templates/pages/login.html",
			},
			// Domain: user
			{
				Source:      "features/auth/github-oauth2/internal/domain/user/model.go.tmpl",
				Destination: "internal/domain/user/model.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/domain/user/repository.go.tmpl",
				Destination: "internal/domain/user/repository.go",
			},
			// Domain: session
			{
				Source:      "features/auth/github-oauth2/internal/domain/session/model.go.tmpl",
				Destination: "internal/domain/session/model.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/domain/session/repository.go.tmpl",
				Destination: "internal/domain/session/repository.go",
			},
			// Application
			{
				Source:      "features/auth/github-oauth2/internal/application/auth/service.go.tmpl",
				Destination: "internal/app/auth/service.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/application/auth/type.go.tmpl",
				Destination: "internal/app/auth/type.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/application/auth/ports.go.tmpl",
				Destination: "internal/app/auth/ports.go",
			},
			// Infrastructure: auth provider
			{
				Source:      "features/auth/github-oauth2/internal/infrastructure/auth/github_oauth.go.tmpl",
				Destination: "internal/infrastructure/auth/github_oauth.go",
			},
			// Infrastructure: persistence (user)
			{
				Source:      "features/auth/github-oauth2/internal/infrastructure/persistence/user_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/user_repository_sqlite.go",
			},
			{
				Source:      "features/auth/github-oauth2/internal/infrastructure/persistence/tx.go.tmpl",
				Destination: "internal/infrastructure/persistence/tx.go",
			},
			// Infrastructure: persistence (session)
			{
				Source:      "features/auth/github-oauth2/internal/infrastructure/persistence/session_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/session_repository_sqlite.go",
			},
			// Infrastructure: http cookie helpers
			{
				Source:      "features/auth/github-oauth2/internal/transport/http/cookies.go.tmpl",
				Destination: "internal/transport/http/cookies.go",
			},
			// Migrations
			{
				Source:      "features/auth/github-oauth2/db/migrations/0001_create_users.sql.tmpl",
				Destination: "db/migrations/0001_create_users.sql",
			},
			{
				Source:      "features/auth/github-oauth2/db/migrations/0002_create_sessions.sql.tmpl",
				Destination: "db/migrations/0002_create_sessions.sql",
			},
			{
				Source:      "features/auth/github-oauth2/db/migrations/0003_create_user_identities.sql.tmpl",
				Destination: "db/migrations/0003_create_user_identities.sql",
			},
		},
	},
}

var defaultSelection = Selection{
	CategoryFrontend: {"frontend-htmx"},
	CategoryStyling:  {"styling-tailwind"},
	CategoryHTTP:     {"http-standard"},
	CategoryDatabase: {"database-none"},
}

// Categories returns a copy of the registered feature categories ordered for display.
func Categories() []FeatureCategory {
	out := make([]FeatureCategory, len(categories))
	copy(out, categories)
	return out
}

// FeaturesForCategory lists features for the supplied category ordered by name.
func FeaturesForCategory(categoryID string) []Feature {
	features := make([]Feature, 0)
	for _, feature := range featureCatalog {
		if feature.CategoryID == categoryID {
			features = append(features, feature)
		}
	}
	sort.Slice(features, func(i, j int) bool {
		return features[i].Name < features[j].Name
	})
	return features
}

// DefaultSelection returns a copy of the default feature selection.
func DefaultSelection() Selection {
	return CloneSelection(defaultSelection)
}

func allFeatures() []Feature {
	out := make([]Feature, len(featureCatalog))
	copy(out, featureCatalog)
	return out
}

// FeatureByID returns the feature for the given identifier.
func FeatureByID(id string) (Feature, bool) {
	for _, feature := range featureCatalog {
		if feature.ID == id {
			return feature, true
		}
	}
	return Feature{}, false
}
