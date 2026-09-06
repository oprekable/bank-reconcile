# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

```bash
# Build & run
make run                    # Build and run app
make run-sample             # Generate sample transaction data
make run-process            # Run reconciliation process

# Testing
make test                   # Run tests with coverage
go test ./path/to/pkg/...   # Run single package tests

# Development workflow
make development-checks     # Full dev cycle: install tools, generate code, lint, staticcheck, deadcode, govulncheck
make generate               # Run wire (DI) and mockery (mocks)

# Linting & static analysis
make go-lint               # golangci-lint
make staticcheck           # staticcheck
make govulncheck           # vulnerability check
make godeadcode            # dead code detection
```

## Architecture

CLI transaction reconciliation system: matches internal system transactions against multiple bank statements, identifies discrepancies.

**Core flow:**
- `main.go` → `cmd/root` (cobra root command)
- Subcommands: `sample` (generate test data), `process` (run reconciliation), `version`
- Each subcommand: `cmd/<name>` → `internal/app/handler/hcli/<name>` → `internal/app/service/<name>` → `internal/app/repository/<name>`

**Dependency injection:** Google Wire (`internal/_inject/wire_gen.go` - auto-generated, don't edit directly). Run `make generate` after changing dependencies.

**Database:** SQLite via `modernc.org/sqlite` (pure Go, CGO_ENABLED=0). Schema and queries in repository layer. Debug mode (`--debug=true`) persists DB to `sample.db` or `reconciliation.db`.

**Configuration:** TOML files in `params/` (embedded via `//go:embed all:embeds`). Config structs in `internal/app/config/`. Environment overrides via `params/.env`.

**Bank parsers:** Strategy pattern in `internal/pkg/reconcile/parser/banks/`. Each bank (BCA, BNI, etc.) has custom CSV format. Default parser for unknown banks. Registry pattern selects parser by bank name.

**Components:** Modular system (`internal/app/component/`) - config, logger, filesystem, sqlite, profiler, error handling. Wired via component sets.

## Key Patterns

- **Interfaces & mocks:** `//go:generate mockery` directives above interface definitions. Mocks in `_mock/` directories.
- **Testing:** `testify` for assertions, `go-sqlmock` for DB mocks, `spf13/afero` for filesystem mocks.
- **Error handling:** Custom error types in `internal/app/err/`. Wrap errors with `pkg/errors`.
- **Logging:** `rs/zerolog` with structured fields. Component-based initialization.
- **Profiling:** pprof support via `--profiler=true` flag. Generates cpu/mem/block/mutex/trace profiles.
- **Progress bars:** `schollz/progressbar/v3` for CLI feedback.

## Code Structure

- `cmd/` - CLI command definitions (cobra)
- `internal/app/` - Application layer (handlers, services, repositories, config, components)
- `internal/pkg/` - Reusable packages (parsers, utils, drivers, shutdown)
- `embeds/` - Embedded config files (params/*.toml)
- `variable/` - Build-time variables (version, git commit, build date)

## Testing

Tests use table-driven patterns. SQL mocks verify query execution. File operations use `afero.MemMapFs` for isolation. Time mocking via interface abstraction.

Run specific test: `go test -v ./internal/app/service/sample/...`

## Build & Release

- **goreleaser** config in `.goreleaser.yaml` - builds for linux/darwin/windows (amd64/arm64)
- Build flags inject version/commit/date via ldflags
- GitHub Actions for CI/CD (`.github/workflows/`)
