# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

gemini-bridge is a protocol bridge that converts Gemini Protocol gemtext content to HTTPS-accessible HTML. It consists of two components:

1. **SSG (Go)** — A static site generator that parses `.gmi` (gemtext) files and produces HTML, metadata JSON, and Atom feeds at build time
2. **Gateway (TypeScript/Hono)** — A Cloudflare Workers application that serves content with proper Gemini-to-HTTP semantic mapping

The design doc is in `gemini-bridge 技術設計書.md` (Japanese).

## Build & Test Commands

```bash
# Build the SSG
go build ./cmd/gemini-bridge/

# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal/domain/parser/
go test ./internal/domain/renderer/

# Run a single test
go test ./internal/domain/parser/ -run TestParseGemtext
```

## Architecture

### SSG Component (Go) — Clean Architecture

```
cmd/gemini-bridge/main.go          — CLI entrypoint, manual DI wiring
internal/
  domain/
    model/       — gemtext AST nodes, document model, front matter, site model
    parser/      — gemtext parser (line-oriented state machine), front matter parser
    renderer/    — HTML renderer using html/template
    feed/        — Atom feed generator
  port/          — Interfaces (ContentWriter, MetadataStore)
  infrastructure/ — Implementations (filesystem writer, JSON metadata store)
  application/   — BuildPipeline orchestration, config
```

### Gateway Component (TypeScript/Hono on Cloudflare Workers)

```
domain/gemini/
  types.ts        — GeminiStatusCode, GeminiResponse, statusCategory()
  semantics.ts    — mapGeminiToHttp() status/header mapping
  negotiation.ts  — Content Negotiation (Accept header → html/gemtext/json)
```

### Key Design Decisions

- **Zero external dependencies for Go SSG** — uses only the standard library (`html/template`, `testing`, etc.)
- **Gemini semantics are first-class** — all HTTP responses include `X-Gemini-Status` and `X-Gemini-Meta` headers preserving original Gemini protocol information
- **Content Negotiation** — `Accept: text/gemini` returns raw gemtext, `text/html` returns SSG-generated HTML, `application/json` returns metadata
- **gemtext AST is flat** (not hierarchical) — nodes are a plain slice, matching gemtext's line-oriented nature
- **Sealed interface pattern** — `Node` interface uses unexported `sealed()` method to prevent external implementations
- **Build-time vs runtime separation** — all parsing/rendering happens in CI/CD; Workers only handle protocol bridging and content delivery

### Gemini Status Code Mapping (core domain logic)

The central design maps all Gemini 2-digit status codes to HTTP equivalents:
- 1x INPUT → 200 + HTML form
- 20 SUCCESS → 200
- 3x REDIRECT → 301/302 + Location header
- 4x TEMPORARY FAILURE → 503/502/429
- 5x PERMANENT FAILURE → 404/410/400
- 6x CLIENT CERTIFICATE → 401/403

### Infrastructure

- **Cloudflare Workers** (gateway runtime), **R2** (static assets), **KV** (cache), **D1** (SQLite metadata), **Workers AI** (summaries)
- **GitHub Actions** for CI/CD (build SSG → deploy to R2/D1 → deploy Workers)
- Gateway tests use **Vitest + Miniflare**

## Git Workflow

- **Branch strategy**: trunk-based development with short-lived branches
  - `main` — always deployable
  - `feat/<name>`, `fix/<name>`, `refactor/<name>`, `docs/<name>`, `ci/<name>`
- **Commit convention**: [Conventional Commits](https://www.conventionalcommits.org/) (English, imperative mood, lowercase, ≤72 chars)
  - Types: `feat`, `fix`, `refactor`, `test`, `docs`, `ci`, `build`, `perf`, `style`, `chore`
  - Scopes: `ssg`, `gateway`, `parser`, `renderer`, `feed`, `infra`
  - Examples: `feat(parser): add front matter extraction`, `fix(renderer): escape ampersands in link text`
- **Commit granularity**: one commit = one logical change; tests belong with the code they cover; refactoring is a separate commit

## Authentication and Signing (1Password)

All configured in global `~/.gitconfig`:
- **SSH authentication**: `core.sshCommand = ssh.exe` delegates to Windows 1Password SSH agent
- **Commit signing**: `gpg.format = ssh`, `gpg.ssh.program = op-ssh-sign-wsl.exe`
- **Signature verification**: `gpg.ssh.allowedSignersFile = ~/.ssh/allowed_signers`
- **GitHub credentials**: `gh auth git-credential`
- **Remote URL format**: `git@github.com:P4suta/<repo>.git` (SSH)

## Security

Never commit these files (enforced by `.gitignore`):
- `.env`, `.env.*` — environment variables / secrets
- `*.pem`, `*.key`, `*.p12` — private keys and certificates
- `wrangler.toml` — may contain account IDs and secrets
- `.dev.vars` — Wrangler local secrets
- `secrets/` — any secret material
