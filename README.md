# Fullkek Starter

Scaffold hypermedia-first Go web apps in seconds.

## Install

```sh
go install github.com/Parapheen/fullkek-starter@latest
```

Or run from source:

```sh
go run . --help
```

## TL;DR (Quickstart)

```sh
# 1) Create a new app using the interactive wizard
fullkek-starter new

# 2) Next steps after generation
make go
```

The server starts and serves a home page at http://localhost:8080 (port can vary by templates). Start editing templates under `web/templates` and assets under `web/assets`.

## What it generates

- `cmd/server`: main entrypoint to start the HTTP server
- `internal/app`: app wiring and composition root
- `internal/runtime`: runtime bootstrap code
- `internal/transport/http`: router and server setup
- `web/templates`: HTML pages and partials
- `web/assets`: CSS/JS aligned with the chosen stack

Stack features add extra files and folders on top of this base.

## Interactive wizard

Just run:

```sh
fullkek new
```

You’ll be guided to enter:

- App name, module path, output directory
- One option for each category: frontend, styling, web framework
- Whether to overwrite the destination if it exists

Press Tab to move forward, Shift+Tab to go back, and Enter to confirm. Esc cancels.

## Non-interactive flags

```sh
fullkek new [app-name]
  --module   string   Go module path (defaults to a sanitized app name)
  --output,-o string  target directory (defaults to app name)
  --force             overwrite destination directory if it exists
  --no-ui             skip the interactive wizard
  --frontend string   frontend feature id
  --styling  string   styling feature id
  --http     string   HTTP framework feature id
  -v, --verbose       verbose output
```

If `--no-ui` is used and `[app-name]` is omitted, the command will error. Without `--no-ui`, leaving `[app-name]` empty opens the wizard.

## Available feature IDs

- Frontend runtime:

  - `frontend-htmx`: HTMX request/response swaps
  - `frontend-fixi`: Fixi.js micro-library bindings

- Styling system:

  - `styling-tailwind`: Tailwind CLI setup with `input.css` and config
  - `styling-daisyui`: DaisyUI standalone bundle delivered over CDN, no build step
  - `styling-tailwind-basecoat`: Tailwind CLI plus Basecoat component library via CDN

- Web framework:

  - `http-standard`: Go net/http ServeMux with HTML response
  - `http-chi`: chi router with starter endpoints

## License

MIT
