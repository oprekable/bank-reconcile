---
name: bank-reconcile-dev
description: Expert development guidelines, architecture patterns, data integrity, and coding rules for the bank-reconcile Go codebase. Use when adding bank parsers, modifying reconciliation algorithms, working with SQLite transactions, updating Google Wire DI, or writing unit tests.
---

# Bank Reconcile Development Guidelines

This skill provides domain-specific knowledge, architectural constraints, and development standards for the `bank-reconcile` CLI application.

---

## 1. Project Overview & Architecture

`bank-reconcile` is a high-performance Go CLI application that reconciles internal system transaction records against external bank statements (BCA, BNI, BSI, Mandiri, BRI, Danamon, etc.) using parallel batch processing and local SQLite storage.

### Layering Convention:

- `cmd/`: CLI layer built on `spf13/cobra`. Handles flags, help docs, and command runners.
- `internal/_inject/`: Dependency injection providers and generated graphs via `google/wire`.
- `internal/app/component/`: Low-level system components (`cconfig`, `clogger`, `csqlite`).
- `internal/app/service/`: Core reconciliation business logic, worker pool partitioning, and validation.
- `internal/app/repository/`: Data layer, queries, and SQLite transaction handlers.
- `internal/pkg/reconcile/parser/`: Pluggable bank parser registry and statement format adapters.
- `internal/pkg/utils/`: File I/O (`spf13/afero`), progress bars, CSV utilities, and logging.

---

## 2. Critical Development Rules & Invariants

### A. Database Transactions & Error Handling

1. **Always use transaction helpers**: Always use `defer helper.CommitOrRollback(tx, &err)` on database transactions with a named `err error` return value.
2. **Never swallow rollback errors**: Ensure error wrapping maintains the original query failure when a rollback succeeds.
3. **Read/Write Connection Separation**: Use `csqlite.DBWrite` exclusively for write/mutation operations (single writer lock) and `csqlite.DBRead` for concurrent read queries.

### B. Concurrency & Partitioning

1. **Worker Partition Calculation**: When dividing amount ranges across workers, calculate chunk steps using:
    ```go
    step = (maxAmount - minAmount) / float64(numWorkers)
    ```
2. **Range Validation**: Always guard against `minAmount > maxAmount`, `minAmount == maxAmount`, and `NumberWorker <= 0` before launching goroutines.
3. **Goroutine Leaks**: Always use `sync.WaitGroup` or `golang.org/x/sync/errgroup` with bounded channels and context cancellation.

### C. Financial Data Integrity

1. **Monetary Precision**: Avoid comparing float values with direct equality (`==`). Prefer tolerance ranges, integer cents (`int64`), or fixed-precision parsing.
2. **CSV Sanitization**: Sanitize numeric and date fields before parsing to handle localized decimal points (`,` vs `.`) and currency symbols.

### D. Filesystem Operations

1. **Use Afero FS**: Never call `os.Open` or `os.RemoveAll` directly in services or parsers. Always inject and use `afero.Fs` (`comp.GetAferoFs()`) to ensure full testability.
2. **Safe Directory Cleanup**: Never delete parent or wildcard directories (`basePath + "/*"`). Delete only specifically generated temporary files.

---

## 3. Extending the Codebase

### Adding a New Bank Parser

1. Create a new package under `internal/pkg/reconcile/parser/banks/<bank_name>/`.
2. Implement the `parser.Parser` interface (`Parse(ctx, reader) ([]model.Transaction, error)`).
3. Register the new parser in `internal/pkg/reconcile/parser/banks/registry.go`.
4. Add comprehensive unit tests with sample bank CSV fixtures.

### Updating Dependency Injection (Google Wire)

Whenever a new service, repository, or component constructor is added:

1. Update `internal/_inject/inject.go`.
2. Run code generation:
    ```bash
    wire ./internal/_inject
    ```

---

## 4. Testing & Validation Guidelines

1. **Table-Driven Tests**: Use table-driven tests (`tests := []struct{...}`) for all parser and service unit tests.
2. **Mocking**: Use `github.com/DATA-DOG/go-sqlmock` for repository tests and `afero.NewMemMapFs()` for filesystem tests.
3. **Running Tests**: Run tests with explicit cache settings:
    ```bash
    GOCACHE=$(pwd)/.gocache go test -v ./...
    ```
4. **Linters & Verification**:
    ```bash
    make go-lint
    make staticcheck
    make govulncheck
    ```
