# Copilot Instructions for zeebe-changelog (zcl)

## Project Overview

This is a Go CLI tool (`zcl`) that generates changelogs for the [Camunda 8](https://github.com/camunda/camunda) project. It interacts with the GitHub API to label issues/PRs and generate markdown-formatted changelogs grouped by component and type.

### Two main commands

- **`add-labels`** — Parses git merge history between two revisions, extracts issue IDs, and adds a GitHub label to each referenced issue/PR. Supports concurrent workers and dry-run mode.
- **`generate`** — Fetches all issues/PRs with a given GitHub label and outputs a structured markdown changelog.

## Tech Stack & Dependencies

- **Language:** Go (module: `github.com/camunda/zeebe-changelog`)
- **CLI framework:** `github.com/urfave/cli` v1 (not v2)
- **GitHub API client:** `github.com/google/go-github/v83`
- **Auth:** `golang.org/x/oauth2` (token-based)
- **Progress bar:** `github.com/gosuri/uiprogress`
- **Testing:** `github.com/stretchr/testify/assert`
- **Vendored dependencies:** All deps are vendored under `vendor/`. Use `-mod=vendor` when building/testing.

## Project Structure

```
cmd/zcl/main.go        — CLI entrypoint, flag definitions, command handlers
pkg/github/client.go   — GitHub API client wrapper (add labels, fetch issues)
pkg/github/changelog.go — Changelog model and markdown rendering
pkg/github/issue.go    — Issue model with label classification helpers
pkg/github/section.go  — Section model (groups issues by component scope)
pkg/gitlog/gitlog.go   — Git log parsing and issue ID extraction
pkg/progress/progress.go — Progress bar wrapper
```

## Build & Test

```sh
make build       # Build binary to bin/zcl
make test        # Run tests with vendored deps
make fmt         # Format Go files with gofmt
make cover       # Run tests with coverage
```

Tests use `-mod=vendor`. Always run `make test` after changes.

## Coding Conventions

- Use `gofmt` for formatting (enforced by Makefile).
- Dependencies are vendored — run `go mod vendor` after adding/updating deps.
- Follow existing patterns: constructors named `NewXxx()`, receiver methods on pointer types.
- Configuration uses CLI flags with environment variable fallbacks (e.g., `--token` / `GITHUB_TOKEN`).
- Constants for flag names, env vars, and labels are defined at the top of each file.
- Error handling: use `log.Fatalln()` for fatal errors, `log.Printf()` for warnings.
- Exported types use unexported fields with accessor methods.
- Test files are colocated with source files (e.g., `changelog_test.go` next to `changelog.go`).
- Tests use `testify/assert` — prefer `assert.Equal`, `assert.True`, etc.

## Domain Concepts

### Issue Labels (in `pkg/github/issue.go`)

Issues are classified by GitHub labels:
- **Scope labels:** `scope/broker`, `scope/gateway`, `scope/clients-java`, `scope/clients-go`, `scope/zbctl`
- **Kind labels:** `kind/feature`, `kind/bug`, `kind/documentation`, `kind/toil`

### Changelog Structure (in `pkg/github/changelog.go`)

The generated changelog has these sections:
1. **Enhancements** — issues with `kind/feature`, grouped by scope (Broker, Gateway, Java Client, Go Client, zbctl, Misc)
2. **Bug Fixes** — issues with `kind/bug`, grouped by scope
3. **Maintenance** — issues with `kind/toil`
4. **Documentation** — issues with `kind/documentation`
5. **Merged Pull Requests** — items that are PRs rather than issues

### Git Log Parsing (in `pkg/gitlog/gitlog.go`)

Extracts issue IDs from merge commit messages by matching patterns like:
- `closes #1234`, `resolves camunda/camunda#5678`
- `merges https://github.com/camunda/camunda/5678`
- Supports `camunda/camunda` repository
