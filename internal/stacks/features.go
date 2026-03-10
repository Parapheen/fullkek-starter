package stacks

import "sort"

const (
	CategoryFrontend       = "frontend"
	CategoryStyling        = "styling"
	CategoryHTTP           = "http"
	CategoryDatabase       = "database"
	CategoryAuth           = "auth"
	CategoryOAuthProviders = "oauth-providers"
	CategoryEmail          = "email"
	CategoryPayments       = "payments"
	CategoryDeploy         = "deploy"
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
	{
		ID:            CategoryOAuthProviders,
		Name:          "OAuth providers",
		Description:   "Select one or more OAuth identity providers.",
		AllowMultiple: true,
	},
	{
		ID:            CategoryEmail,
		Name:          "Email",
		Description:   "Email sending strategy.",
		AllowMultiple: false,
	},
	{
		ID:            CategoryPayments,
		Name:          "Payments",
		Description:   "Payment processing integration.",
		AllowMultiple: false,
	},
	{
		ID:            CategoryDeploy,
		Name:          "Deployment",
		Description:   "Deployment automation for your project.",
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
				Source:      "features/frontend/htmx/assets/scripts/htmx.min.js.tmpl",
				Destination: "public/assets/scripts/htmx.min.js",
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
	// --- Authentication ---
	{
		ID:          "auth-none",
		CategoryID:  CategoryAuth,
		Name:        "None",
		Description: "Skip bundling authentication logic.",
		Tags:        []string{"auth"},
	},
	{
		ID:          "auth-oauth2",
		CategoryID:  CategoryAuth,
		Name:        "OAuth2",
		Description: "Login with OAuth2 providers using server-side sessions.",
		Tags:        []string{"auth", "oauth2"},
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
				Source:      "features/auth/oauth2/internal/transport/http/oauth_handlers.go.tmpl",
				Destination: "internal/transport/http/oauth_handlers.go",
			},
			{
				Source:      "features/auth/oauth2/internal/transport/http/render.go.tmpl",
				Destination: "internal/transport/http/render.go",
			},
			{
				Source:      "features/auth/oauth2/internal/transport/http/auth_middleware.go.tmpl",
				Destination: "internal/transport/http/auth_middleware.go",
			},
			// Web templates
			{
				Source:      "features/auth/oauth2/web/templates/pages/profile.html.tmpl",
				Destination: "web/templates/pages/profile.html",
			},
			{
				Source:      "features/auth/oauth2/web/templates/pages/login.html.tmpl",
				Destination: "web/templates/pages/login.html",
			},
			// Domain: user
			{
				Source:      "features/auth/oauth2/internal/domain/user/model.go.tmpl",
				Destination: "internal/domain/user/model.go",
			},
			{
				Source:      "features/auth/oauth2/internal/domain/user/repository.go.tmpl",
				Destination: "internal/domain/user/repository.go",
			},
			// Domain: session
			{
				Source:      "features/auth/oauth2/internal/domain/session/model.go.tmpl",
				Destination: "internal/domain/session/model.go",
			},
			{
				Source:      "features/auth/oauth2/internal/domain/session/repository.go.tmpl",
				Destination: "internal/domain/session/repository.go",
			},
			// Application
			{
				Source:      "features/auth/oauth2/internal/application/auth/service.go.tmpl",
				Destination: "internal/app/auth/service.go",
			},
			{
				Source:      "features/auth/oauth2/internal/application/auth/type.go.tmpl",
				Destination: "internal/app/auth/type.go",
			},
			{
				Source:      "features/auth/oauth2/internal/application/auth/ports.go.tmpl",
				Destination: "internal/app/auth/ports.go",
			},
			// Infrastructure: persistence (user)
			{
				Source:      "features/auth/oauth2/internal/infrastructure/persistence/user_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/user_repository_sqlite.go",
			},
			{
				Source:      "features/auth/oauth2/internal/infrastructure/persistence/tx.go.tmpl",
				Destination: "internal/infrastructure/persistence/tx.go",
			},
			// Infrastructure: persistence (session)
			{
				Source:      "features/auth/oauth2/internal/infrastructure/persistence/session_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/session_repository_sqlite.go",
			},
			// Infrastructure: http cookie helpers
			{
				Source:      "features/auth/oauth2/internal/transport/http/cookies.go.tmpl",
				Destination: "internal/transport/http/cookies.go",
			},
			// Migrations
			{
				Source:      "features/auth/oauth2/db/migrations/0001_create_users.sql.tmpl",
				Destination: "db/migrations/0001_create_users.sql",
			},
			{
				Source:      "features/auth/oauth2/db/migrations/0002_create_sessions.sql.tmpl",
				Destination: "db/migrations/0002_create_sessions.sql",
			},
			{
				Source:      "features/auth/oauth2/db/migrations/0003_create_user_identities.sql.tmpl",
				Destination: "db/migrations/0003_create_user_identities.sql",
			},
		},
	},
	{
		ID:          "auth-magic-link",
		CategoryID:  CategoryAuth,
		Name:        "Magic Link",
		Description: "Passwordless sign-in via one-time emailed link (logged in development).",
		Tags:        []string{"auth", "magic-link", "passwordless"},
		Directories: []string{
			"db/migrations",
			"internal/app/auth",
			"internal/domain/magiclink",
			"internal/domain/session",
			"internal/domain/user",
			"internal/infrastructure/persistence",
			"internal/transport/http",
			"web/templates/pages",
		},
		Templates: []Template{
			// HTTP transport
			{
				Source:      "features/auth/magic-link/internal/transport/http/auth_handlers.go.tmpl",
				Destination: "internal/transport/http/auth_handlers.go",
			},
			{
				Source:      "features/auth/magic-link/internal/transport/http/render.go.tmpl",
				Destination: "internal/transport/http/render.go",
			},
			{
				Source:      "features/auth/magic-link/internal/transport/http/auth_middleware.go.tmpl",
				Destination: "internal/transport/http/auth_middleware.go",
			},
			{
				Source:      "features/auth/magic-link/internal/transport/http/cookies.go.tmpl",
				Destination: "internal/transport/http/cookies.go",
			},
			// Web templates
			{
				Source:      "features/auth/magic-link/web/templates/pages/profile.html.tmpl",
				Destination: "web/templates/pages/profile.html",
			},
			{
				Source:      "features/auth/magic-link/web/templates/pages/login.html.tmpl",
				Destination: "web/templates/pages/login.html",
			},
			// Domain: user
			{
				Source:      "features/auth/magic-link/internal/domain/user/model.go.tmpl",
				Destination: "internal/domain/user/model.go",
			},
			{
				Source:      "features/auth/magic-link/internal/domain/user/repository.go.tmpl",
				Destination: "internal/domain/user/repository.go",
			},
			// Domain: session
			{
				Source:      "features/auth/magic-link/internal/domain/session/model.go.tmpl",
				Destination: "internal/domain/session/model.go",
			},
			{
				Source:      "features/auth/magic-link/internal/domain/session/repository.go.tmpl",
				Destination: "internal/domain/session/repository.go",
			},
			// Domain: magic link
			{
				Source:      "features/auth/magic-link/internal/domain/magiclink/model.go.tmpl",
				Destination: "internal/domain/magiclink/model.go",
			},
			{
				Source:      "features/auth/magic-link/internal/domain/magiclink/repository.go.tmpl",
				Destination: "internal/domain/magiclink/repository.go",
			},
			// Application
			{
				Source:      "features/auth/magic-link/internal/application/auth/service.go.tmpl",
				Destination: "internal/app/auth/service.go",
			},
			{
				Source:      "features/auth/magic-link/internal/application/auth/type.go.tmpl",
				Destination: "internal/app/auth/type.go",
			},
			{
				Source:      "features/auth/magic-link/internal/application/auth/ports.go.tmpl",
				Destination: "internal/app/auth/ports.go",
			},
			// Infrastructure: persistence
			{
				Source:      "features/auth/magic-link/internal/infrastructure/persistence/user_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/user_repository_sqlite.go",
			},
			{
				Source:      "features/auth/magic-link/internal/infrastructure/persistence/session_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/session_repository_sqlite.go",
			},
			{
				Source:      "features/auth/magic-link/internal/infrastructure/persistence/magic_link_token_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/magic_link_token_repository_sqlite.go",
			},
			{
				Source:      "features/auth/magic-link/internal/infrastructure/persistence/tx.go.tmpl",
				Destination: "internal/infrastructure/persistence/tx.go",
			},
			// Migrations
			{
				Source:      "features/auth/magic-link/db/migrations/0001_create_users.sql.tmpl",
				Destination: "db/migrations/0001_create_users.sql",
			},
			{
				Source:      "features/auth/magic-link/db/migrations/0002_create_sessions.sql.tmpl",
				Destination: "db/migrations/0002_create_sessions.sql",
			},
			{
				Source:      "features/auth/magic-link/db/migrations/0003_create_magic_link_tokens.sql.tmpl",
				Destination: "db/migrations/0003_create_magic_link_tokens.sql",
			},
		},
	},
	// --- OAuth Providers ---
	{
		ID:          "oauth-github",
		CategoryID:  CategoryOAuthProviders,
		Name:        "GitHub",
		Description: "GitHub OAuth2 identity provider.",
		Tags:        []string{"github"},
		Directories: []string{
			"internal/infrastructure/auth",
		},
		Templates: []Template{
			{
				Source:      "features/oauth/github/internal/infrastructure/auth/github_oauth.go.tmpl",
				Destination: "internal/infrastructure/auth/github_oauth.go",
			},
		},
	},
	{
		ID:          "oauth-google",
		CategoryID:  CategoryOAuthProviders,
		Name:        "Google",
		Description: "Google OAuth2 identity provider.",
		Tags:        []string{"google"},
		Directories: []string{
			"internal/infrastructure/auth",
		},
		Templates: []Template{
			{
				Source:      "features/oauth/google/internal/infrastructure/auth/google_oauth.go.tmpl",
				Destination: "internal/infrastructure/auth/google_oauth.go",
			},
		},
	},
	{
		ID:          "oauth-yandex",
		CategoryID:  CategoryOAuthProviders,
		Name:        "Yandex",
		Description: "Yandex OAuth2 identity provider.",
		Tags:        []string{"yandex"},
		Directories: []string{
			"internal/infrastructure/auth",
		},
		Templates: []Template{
			{
				Source:      "features/oauth/yandex/internal/infrastructure/auth/yandex_oauth.go.tmpl",
				Destination: "internal/infrastructure/auth/yandex_oauth.go",
			},
		},
	},
	// --- Email ---
	{
		ID:          "email-none",
		CategoryID:  CategoryEmail,
		Name:        "None",
		Description: "Skip email (logs to stdout in dev mode).",
		Tags:        []string{"email"},
	},
	{
		ID:          "email-smtp",
		CategoryID:  CategoryEmail,
		Name:        "SMTP",
		Description: "Send email via SMTP with STARTTLS support.",
		Tags:        []string{"email", "smtp"},
		Directories: []string{
			"internal/app/email",
			"internal/infrastructure/email",
		},
		Templates: []Template{
			{
				Source:      "features/email/smtp/internal/app/email/ports.go.tmpl",
				Destination: "internal/app/email/ports.go",
			},
			{
				Source:      "features/email/smtp/internal/infrastructure/email/smtp.go.tmpl",
				Destination: "internal/infrastructure/email/smtp.go",
			},
		},
	},
	// --- Payments ---
	{
		ID:          "payments-none",
		CategoryID:  CategoryPayments,
		Name:        "None",
		Description: "Skip payment integration.",
		Tags:        []string{"payments"},
	},
	{
		ID:          "payments-yookassa",
		CategoryID:  CategoryPayments,
		Name:        "YooKassa",
		Description: "YooKassa checkout integration.",
		Tags:        []string{"payments", "yookassa"},
		Directories: []string{
			"db/migrations",
			"internal/domain/payment",
			"internal/infrastructure/payments",
			"internal/infrastructure/persistence",
			"internal/transport/http",
			"web/templates/pages",
		},
		Templates: []Template{
			{
				Source:      "features/payments/yookassa/internal/domain/payment/model.go.tmpl",
				Destination: "internal/domain/payment/model.go",
			},
			{
				Source:      "features/payments/yookassa/internal/domain/payment/repository.go.tmpl",
				Destination: "internal/domain/payment/repository.go",
			},
			{
				Source:      "features/payments/yookassa/internal/infrastructure/payments/yookassa.go.tmpl",
				Destination: "internal/infrastructure/payments/yookassa.go",
			},
			{
				Source:      "features/payments/yookassa/internal/infrastructure/persistence/payment_repository_sqlite.go.tmpl",
				Destination: "internal/infrastructure/persistence/payment_repository_sqlite.go",
			},
			{
				Source:      "features/payments/yookassa/internal/transport/http/payment_handlers.go.tmpl",
				Destination: "internal/transport/http/payment_handlers.go",
			},
			{
				Source:      "features/payments/yookassa/web/templates/pages/checkout.html.tmpl",
				Destination: "web/templates/pages/checkout.html",
			},
			{
				Source:      "features/payments/yookassa/web/templates/pages/payment_success.html.tmpl",
				Destination: "web/templates/pages/payment_success.html",
			},
			{
				Source:      "features/payments/yookassa/db/migrations/0004_create_payments.sql.tmpl",
				Destination: "db/migrations/0004_create_payments.sql",
			},
		},
	},
	// --- Deployment ---
	{
		ID:          "deploy-none",
		CategoryID:  CategoryDeploy,
		Name:        "None",
		Description: "Skip deployment automation.",
		Tags:        []string{"deploy"},
	},
	{
		ID:          "deploy-ansible",
		CategoryID:  CategoryDeploy,
		Name:        "Ansible",
		Description: "Ansible playbook for Ubuntu VPS with Caddy reverse proxy.",
		Tags:        []string{"deploy", "ansible"},
		Directories: []string{
			"deploy/templates",
			"deploy/group_vars",
		},
		Templates: []Template{
			{
				Source:      "features/deploy/ansible/deploy/inventory.ini.tmpl",
				Destination: "deploy/inventory.ini",
			},
			{
				Source:      "features/deploy/ansible/deploy/playbook.yml.tmpl",
				Destination: "deploy/playbook.yml",
			},
			{
				Source:      "features/deploy/ansible/deploy/templates/app.service.j2.tmpl",
				Destination: "deploy/templates/app.service.j2",
			},
			{
				Source:      "features/deploy/ansible/deploy/templates/Caddyfile.j2.tmpl",
				Destination: "deploy/templates/Caddyfile.j2",
			},
			{
				Source:      "features/deploy/ansible/deploy/group_vars/all.yml.tmpl",
				Destination: "deploy/group_vars/all.yml",
			},
		},
	},
}

var defaultSelection = Selection{
	CategoryFrontend:       {"frontend-htmx"},
	CategoryStyling:        {"styling-tailwind"},
	CategoryHTTP:           {"http-standard"},
	CategoryDatabase:       {"database-none"},
	CategoryAuth:           {"auth-none"},
	CategoryOAuthProviders: {},
	CategoryEmail:          {"email-none"},
	CategoryPayments:       {"payments-none"},
	CategoryDeploy:         {"deploy-none"},
}

var featureDependencies = map[string][]string{
	"auth-oauth2":       {"database-sqlite"},
	"auth-magic-link":   {"database-sqlite"},
	"oauth-github":      {"auth-oauth2"},
	"oauth-google":      {"auth-oauth2"},
	"oauth-yandex":      {"auth-oauth2"},
	"payments-yookassa": {"database-sqlite"},
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

// FeatureDependencies returns the IDs required by a feature.
func FeatureDependencies(id string) []string {
	deps := featureDependencies[id]
	out := make([]string, len(deps))
	copy(out, deps)
	return out
}
