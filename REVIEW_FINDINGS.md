# Code Review Findings

**Date:** 2026-08-11
**Dimensions reviewed:** Best practices, Memory & CPU, Concurrency, Performance, Security, GC pressure, Maintainability
**Results:** 83 findings generated, 60 confirmed after adversarial verification

---

## HIGH Severity

### 1. Race Condition: Concurrent map write without mutex
**File:** `internal/app/service/process/svc.go:450`
**Category:** Concurrency

`parallel.ForEach(bankNames, ...)` writes to `reconciliationSummary.FileMissingBankTrx` map from multiple goroutines without mutex. Go maps panic on concurrent write.

```go
// BROKEN:
parallel.ForEach(bankNames, func(item string, _ int) {
    reconciliationSummary.FileMissingBankTrx[item] = fileReportBankTrx  // RACE
})
```

**Fix:** Add `sync.Mutex` around map write, or collect results via channel.

---

### 2. Race Condition: `stmtMap` accessed without mutex
**Files:**
- `internal/app/repository/helper/helper.go:30-32`
- `internal/app/repository/process/db.go:24`
- `internal/app/repository/sample/db.go:33`

**Category:** Concurrency

`stmtMap map[string]*sql.Stmt` read/written concurrently in `QueryContext` and `ExecTxQueries` without synchronization. Multiple goroutines calling repository methods concurrently causes panic.

**Fix:** Add `sync.RWMutex` to DB struct, protect all `stmtMap` access. Or use `sync.Map`.

---

### 3. SIGSEGV caught in signal handler
**File:** `internal/pkg/shutdown/shutdown.go:19`
**Category:** Signal handling / Security

```go
signal.Notify(trap, syscall.SIGINT, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSEGV)
```

SIGSEGV = memory corruption signal. Go runtime uses it internally for stack growth and nil checks. Catching it = undefined behavior.

**Fix:** Remove `syscall.SIGSEGV` from signal list.

---

### 4. sync.Pool use-after-Put (data corruption + zero benefit)
**Files:**
- `internal/pkg/reconcile/parser/systems/default_system/system_parser_default.go:121-131`
- `internal/pkg/reconcile/parser/banks/helper/helper.go:54-65`
- `internal/app/service/sample/svc.go:190-205`

**Category:** GC / Concurrency / Correctness

Three compounding bugs:
1. Pool `Get()` result immediately overwritten by fresh allocation — pooled object DISCARDED, pool provides zero reuse benefit
2. Object returned to pool via `Put()`, then mutated after — race condition
3. Caller holds reference to pooled object in return slice

```go
// BROKEN pattern (all 3 locations):
ptr := pool.Get().(*Type)
ptr, err = newValue()     // pooled object DISCARDED
pool.Put(ptr)             // returned to pool
ptr.Field = x             // MUTATED after Put — race
returnData = append(returnData, ptr)  // caller holds pooled ref
```

Note: `//nolint:all` and `//lint:ignore SA4006` comments confirm author saw warnings but suppressed them incorrectly.

**Fix:** Remove all `sync.Pool` usage entirely. Use local variables. Or refactor `ToSystemTrxData()`/`ToBankTrxData()` to populate existing struct + copy before `Put()`.

---

### 5. Swallowed errors in parallel file parsing
**Files:**
- `internal/app/service/process/svc.go:108` — system transaction files
- `internal/app/service/process/svc.go:245` — bank transaction files

**Category:** Error handling / Correctness

```go
data, _ := s.parseSystemTrxFile(ctx, afs, item)  // error DISCARDED
data, _ := s.parseBankTrxFile(c, afs, item)      // error DISCARDED
```

Failed parses silently skipped. Reconciliation runs on incomplete data with no failure signal to operator. Critical for financial data correctness — can produce false-positive "all balanced" when files silently fail.

**Fix:** Collect errors, propagate to caller. Fail or warn when files are skipped.

---

### 6. Arbitrary directory deletion via path flags
**File:** `internal/pkg/utils/csvhelper/csvhelper.go:37`
**Category:** Security / Path traversal

`DeleteDirectory` calls `RemoveAll` on `filepath.Dir(filePath)/*` without base-path validation. Setting `--system-trx-path=/` wipes filesystem root.

```go
func DeleteDirectory(ctx context.Context, fs afero.Fs, filePath string) (err error) {
    basePath := filepath.Dir(filePath)
    afero.Glob(fs, basePath+"/*")
    fs.RemoveAll(item)  // NO VALIDATION
}
```

**Fix:** Validate paths are under expected base directories. Reject `..`, symlinks outside bounds, absolute paths to sensitive locations.

---

## MEDIUM Severity

### 7. Fragile error string comparison
**File:** `internal/app/appcontext/appcontext.go:86`
**Category:** Maintainability

```go
if context.Cause(a.ctx).Error() == "done" { err = nil }
```

Couples to string `"done"` from `cli.go:68`. Renaming either silently breaks CLI success detection. Nil-dereference risk if `Cause` returns nil.

**Fix:** Define sentinel error `var ErrDone = errors.New("done")`, use `errors.Is()`.

---

### 8. Ignored errors across multiple files

| File | Line | Issue |
|------|------|-------|
| `internal/app/repository/process/db.go` | 140 | `b, _ := json.Marshal(data)` — silent data corruption on marshal failure |
| `internal/app/repository/process/db.go` | 101 | `_ = json.NewEncoder(b).Encode(listBank)` — encoding error discarded |
| `internal/pkg/reconcile/parser/systems/default_system/system_parser_default.go` | 104 | `header, _ := csvutil.Header(...)` — invalid headers on failure |
| `internal/app/component/cprofiler/cprofiler.go` | 43 | `dir, _ := os.Getwd()` — empty path on failure |

**Fix:** Check errors, log or return them.

---

### 9. Missing SQLite pragmas for bulk write performance
**File:** `internal/pkg/driver/sql/sqlite.go:26`
**Category:** Performance

DSN only sets `cache` and `journal_mode`. Missing critical pragmas for bulk write:
- `synchronous=OFF` (or NORMAL for WAL)
- `temp_store=MEMORY`
- `cache_size=-20000` (20MB)
- `mmap_size=268435456`
- `busy_timeout=5000`

Expected 10-100x write performance improvement for bulk operations.

Also: `journal_mode=WAL` has no effect on `:memory:` databases.

**Related:** No `sql.DB` connection pool tuning — missing `SetMaxOpenConns(1)` for SQLite single-writer model, `SetMaxIdleConns(2)`, `SetConnMaxLifetime(0)`.

---

### 10. Sequential chunked DB imports defeat NumberWorker purpose
**File:** `internal/app/service/process/svc.go:125` (also `:150`, `:261`)
**Category:** Performance / Concurrency

Data split into `NumberWorker` chunks but processed sequentially in `for` loop. Each chunk opens separate transaction. Codebase already has `parallel.ForEach` infrastructure.

**Fix:** Use `parallel.ForEach` for chunk processing, or merge into single transaction.

---

### 11. Full CSV marshal to memory before file write
**File:** `internal/pkg/utils/csvhelper/csvhelper.go:26`
**Category:** Memory / GC

`csvutil.Marshal` serializes entire struct slice to `[]byte` in memory, then `WriteFile` writes it. Doubles peak memory for large datasets (structs + CSV bytes).

**Fix:** Use `csvutil.NewEncoder` with `bufio.Writer` for streaming rows directly to file. Reduces peak memory by ~50%.

---

### 12. Prepared statements never closed — resource leak
**File:** `internal/app/repository/helper/helper.go:32`
**Category:** Resource leak

`stmtMap` entries created via `db.PrepareContext` never call `stmt.Close()`. `ExecTxQueries` (line 70) deletes map entries without closing statements first. Bounded leak (finite query name set) but wasteful.

Note: `//nolint:sqlclosecheck` directives confirm deliberate suppression.

**Fix:** Add `stmt.Close()` before `delete` in `ExecTxQueries`. Iterate `stmtMap` and close all entries in `DB.Close()`.

---

### 13. Silent nil-dependency guards return success
**Files:**
- `internal/app/handler/hcli/sample/handler.go:38`
- `internal/app/handler/hcli/process/handler.go:50`

**Category:** Error handling

```go
if h.comp == nil || h.svc == nil || h.repo == nil { return nil }
```

Returns `nil` error when dependencies missing. Caller treats nil as success. Should fail loud.

**Fix:** Return descriptive error indicating which dependency is nil.

---

### 14. Full JSON serialization of large data slices for SQL insert
**File:** `internal/app/repository/process/db.go:140`
**Category:** Memory / Performance

Entire data slice serialized to JSON → string copy → SQLite `json_each` re-parses. For thousands of transactions, unnecessary allocations. Data chunked by `NumberWorker` but still wasteful per-chunk.

**Fix:** Batch INSERT with bound parameters instead of JSON round-trip.

---

### 15. Bank parser code duplication
**Files:**
- `internal/pkg/reconcile/parser/banks/bca/bank_parser_bca.go`
- `internal/pkg/reconcile/parser/banks/bni/bank_parser_bni.go`
- `internal/pkg/reconcile/parser/banks/default_bank/bank_parser_default.go`

**Category:** Maintainability

Three parser implementations share identical struct layouts, method signatures, and implementations. Only differ in parser constant and entity type imported. Entity types also duplicated with only field name prefixes (BCA-, BNI-, Default-).

**Fix:** Consolidate into single generic parser with configuration or factory pattern.

---

### 16. CLI handler struct/method duplication
**Files:**
- `internal/app/handler/hcli/process/handler.go:39`
- `internal/app/handler/hcli/sample/handler.go`

**Category:** Maintainability

Process and sample handlers duplicate identical struct fields (`comp`/`svc`/`repo`/`writer`) and setter methods (`SetComponents`/`SetServices`/`SetRepositories`). 27+ lines of character-for-character duplicate code.

**Fix:** Extract `baseHandler` embedded struct with shared fields and setters.

---

### 17. No validation on filesystem path flags
**File:** `cmd/helper/helper.go:31`
**Category:** Security

Path flags (`--report-trx-path`, `--system-trx-path`, `--bank-trx-path`) accept any string without validation. No checks for: paths outside expected directories, symlink targets, null bytes, excessive length, device files.

**Fix:** Validate against allowlist of base directories. Reject suspicious patterns.

---

## LOW Severity

| # | File:Line | Category | Issue |
|---|-----------|----------|-------|
| 18 | `main.go:44` | Anti-pattern | `unsafe.Pointer` for bool→int conversion — use simple `if` |
| 19 | `internal/app/component/csqlite/csqlite.go:36` | Dead code | `sync.Once` on freshly-instantiated struct fields — meaningless, always executes once |
| 20 | `internal/app/server/server.go:24` | Over-engineering | Reflection (`reflect.ValueOf`) to iterate struct fields — use `[]IServer` slice |
| 21 | `internal/app/handler/hcli/register.go:12` | Style | Package-level global state (intentional Wire DI pattern, acceptable) |
| 22 | `internal/pkg/utils/csvhelper/csvhelper.go:13` | Over-engineering | `hunch.Waterfall` for 4 sequential operations — plain function calls clearer |
| 23 | Bank parsers (BCA/BNI/default) `:28` | Naming | Error message mentions `dataStruct` parameter that doesn't exist — copy-paste artifact |
| 24 | `internal/pkg/driver/sql/sqlite.go:26` | Security | DSN injection via unsanitized string concat (low risk for CLI tool, but defense in depth) |
| 25 | `internal/app/service/process/svc.go:103` | Dead code | Redundant `sync.WaitGroup` inside `parallel.ForEach` (ForEach already blocks) |
| 26 | `internal/app/service/sample/svc.go:199` | Robustness | Unsafe type assertion without comma-ok — works today, fragile if new implementations added |
| 27 | `cmd/process/process.go:100` | Readability | Entire `Runner` body inside `if-else` — use early-return guard clause |

---

## Priority Fix Order

1. **Race conditions** (#1, #2) — crash-class defects, fix immediately
2. **sync.Pool use-after-Put** (#4) — correctness bug + zero benefit, remove all pools
3. **SIGSEGV handler** (#3) — remove `syscall.SIGSEGV`
4. **Swallowed parse errors** (#5) — financial data integrity risk
5. **Path traversal** (#6) — add base-path validation
6. **SQLite pragmas + pool tuning** (#9) — 10-100x bulk write improvement
7. **Error string comparison** (#7) — sentinel error refactor
8. **Remaining error handling** (#8) — check and propagate all errors

---

## Verified False Positives / Non-Issues

- `parallel.ForEach` + `WaitGroup` race condition claim — incorrect, `parallel.ForEach` blocks internally, no race exists
- Package-level globals in `register.go` — intentional Wire DI provider pattern, not a defect
