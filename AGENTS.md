# AGENTS.md
Guide for coding agents working in this repository.

## Repository shape
- This repo has two Go modules:
  - Root module: `github.com/Parapheen/fullkek-starter`
  - Nested example app: `fullkek/` (module path `fullkek`)
- Always run commands from the module you are changing.
- `go test ./...` at repo root does not include the nested `fullkek/` module.

## Discovered AI instruction files
- Found Cursor rules at `.cursor/rules.md`.
- No `.cursor/rules/` directory found.
- No `.cursorrules` file found.
- No Copilot instruction file found at `.github/copilot-instructions.md`.

## Cursor rules to honor (distilled)
- Complete work end-to-end; if blocked, state done vs remaining work clearly.
- Read relevant files before editing.
- Prefer `apply_patch` and keep patches semantically focused.
- Run independent read/search operations in parallel when useful.
- Keep responses concise and conclusion-first.
- For high-risk work (security/auth/infra/cost/production), present plan + risks first.
- Ask clarification only when ambiguity materially changes the result.

## Build, run, lint, test commands

### Root module (`/`)
- Install/sync deps: `go mod tidy`
- Build all packages: `go build ./...`
- Build CLI binary: `go build -o bin/fullkek-starter .`
- Run CLI help: `go run . --help`
- Run scaffold command: `go run . new --no-ui my-app`
- Test all packages: `go test ./...`
- Vet: `go vet ./...`
- Format check: `gofmt -l .`

### Nested app module (`/fullkek`)
- Install/sync deps: `go mod tidy`
- Start dev server: `make dev`
- Start watcher + migrations + server: `make go`
- Build app: `make build`
- Build and run app: `make run`
- Direct build: `go build -o bin/fullkek ./cmd/server`
- Test all packages: `make test` or `go test -v ./...`
- Vet: `go vet ./...`
- Format check: `gofmt -l .`

## Running a single test (important)
- Single test function:
  - `go test ./internal/stacks -run '^TestCompose$' -v`
- Single subtest:
  - `go test ./internal/stacks -run 'TestCompose/invalid_feature' -v`
- Multiple tests by regex:
  - `go test ./internal/stacks -run 'TestCompose|TestValidateSelection' -v`
- Same pattern in nested module:
  - `go test ./internal/app/auth -run '^TestService_HandleCallback$' -v`
- Disable test cache while debugging:
  - `go test ./internal/stacks -run '^TestCompose$' -count=1 -v`

## Code style conventions for this codebase

### Imports
- Use normal Go import grouping:
  1) standard library
  2) third-party packages
  3) same-module imports
- Keep imports `gofmt`-clean; do not manually align columns.
- Alias imports only when it improves clarity or avoids collisions.
- Existing alias patterns: `appauth`, `oauthinfra`, `httptransport`, `domainUser`.

### Formatting and structure
- `gofmt` is required.
- Prefer small cohesive functions and helper extraction for repeated logic.
- Prefer early returns over deep nesting.
- Keep comments for non-obvious intent; do not restate code literally.
- Keep files focused around one responsibility (router, service, persistence, etc.).

### Types and interfaces
- Start with concrete types; introduce interfaces at real boundaries.
- Define interfaces where consumed (ports/services), not prematurely.
- Use typed config structs for constructors and bootstrap wiring.
- Keep structs explicit; avoid hidden magic defaults beyond constructor logic.
- Preserve zero-value safety where practical.

### Naming
- Exported identifiers: `PascalCase`.
- Unexported identifiers: `camelCase`.
- Package names: short, lowercase, no underscores.
- File names: lowercase; underscores are acceptable for clarity.
- Existing file naming style includes `oauth_handlers.go`, `user_repository_sqlite.go`.
- Keep IDs and constants descriptive (`CategoryFrontend`, `sqliteFeatureID`, etc.).

### Error handling
- Return errors, do not panic in library/internal code.
- Wrap propagated errors with context using `%w`.
- Prefer actionable error text with operation context.
- Keep error strings lowercase and concise.
- Use `errors.Is`/`errors.As` for typed/sentinel checks.
- At entrypoints (`main`), log fatal errors and exit non-zero.

### Context and side effects
- Pass `context.Context` as first arg for I/O, DB, network, or cancellable operations.
- Check `ctx.Err()` in loops or multi-step operations that can be cancelled.
- Do not store context in struct fields.
- Keep side effects near composition roots (`cmd/`, `internal/app/`).

### HTTP and transport
- Keep handlers thin: parse input, call service, write response.
- Set explicit response `Content-Type`.
- Keep route registration centralized in `router.go`.
- Put cross-cutting concerns in middleware (logging, recovery, CSRF, auth).

### Logging
- Use structured logs (`log/slog`) for runtime paths.
- Keep log keys consistent (for example `"err"`, `"addr"`).
- Never log secrets, OAuth tokens, session IDs, or credentials.

### Persistence/DB
- Keep DB concerns in `internal/infrastructure/persistence`.
- Wrap transaction/repository errors with operation context.
- Follow existing SQLite defaults (busy timeout, FK on, WAL) unless requested otherwise.

## Architecture boundaries to preserve
- `cmd/`: CLI and process entrypoints.
- `internal/app/`: application composition and use-case orchestration.
- `internal/domain/`: domain models and repository contracts.
- `internal/infrastructure/`: external adapters (DB, OAuth, etc.).
- `internal/transport/http/`: HTTP delivery, middleware, routing.
- `internal/scaffold/` and `internal/stacks/`: generator and stack composition logic.

## Safe-change checklist for agents
- Make minimal diffs targeted to the user request.
- Avoid broad renames/moves unless required.
- Keep generated scaffold behavior backward compatible unless asked to change it.
- Update docs/examples when commands or behavior change.
- After edits in each touched module, run:
  - `go test ./...`
  - `go vet ./...`
  - `gofmt -l .` (and `gofmt -w` on changed files)
- For CLI changes, smoke test with `go run . --help` or target command.
- For nested app runtime changes, smoke test with `make dev` or focused package tests.
