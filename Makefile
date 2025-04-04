ifeq ($(OS), Windows_NT)
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
else
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
endif

.DEFAULT_GOAL := run

.PHONY: download
download:
	@echo Download go.mod dependencies
	@go mod download

.PHONY: install-tools
install-tools: download
	@#go install $(go list -e -f '{{join .Imports " "}}' tools.go)
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

.PHONY: generate
generate:
	@wire github.com/oprekable/bank-reconcile/internal/inject
	@go generate ./...

.PHONY: go-lint
go-lint:
	@golangci-lint linters
	@golangci-lint run ./...

.PHONY: staticcheck
staticcheck:
	@staticcheck ./...

.PHONY: govulncheck
govulncheck:
	@govulncheck ./...

.PHONY: godeadcode
godeadcode:
	@deadcode ./...

.PHONY: go-lint-fix-struct-staticcheck-govulncheck
go-lint-fix-struct-staticcheck-govulncheck: install-tools generate
	@go mod tidy
	@golangci-lint run ./... --fix
	@staticcheck ./...
	@fieldalignment -fix ./...
	@deadcode ./...
	@govulncheck -show verbose ./...

.PHONY: test
test:
	@go test -gcflags=all=-l -count=1 -p=8 -parallel=8 -race -coverprofile=coverage.out ./... -json | tee report.json
	@go tool cover -html=coverage.out

.PHONY: run
run:
	@go build -gcflags -live .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile
	@#env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go


UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
	base_args="--showlog=true --listbank=bca,bni,mandiri,bri,danamon --from=$$(date -j -v -10d '+%Y-%m-%d') --to=$$(date -j '+%Y-%m-%d')"
endif

ifeq ($(UNAME), Linux)
	base_args="--showlog=true --listbank=bca,bni,mandiri,bri,danamon --from=$$(date -d '-10 day' '+%Y-%m-%d') --to=$$(date '+%Y-%m-%d')"
endif

process_args="process ${base_args} -i=true -g=false -s=/tmp/sample/system -b=/tmp/sample/bank -r=/tmp/report"
#process_args="process ${base_args} -g=true"
#process_args="process ${base_args} -g=false"
sample_args="sample ${base_args} -i=true --percentagematch=100 --amountdata=100000 -g=true -s=/tmp/sample/system -b=/tmp/sample/bank"
#sample_args="sample ${base_args} --percentagematch=100 --amountdata=1000 -g=true"
#sample_args="sample ${base_args} --percentagematch=100 --amountdata=1000 -g=false"

.PHONY: echo-sample-args
echo-sample-args:
	@echo $(sample_args)
	@echo $(sample_args) | pbcopy

.PHONY: run-sample
run-sample:
	@echo $(sample_args)
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile  $$(echo $(sample_args))
	@#env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go  $$(echo $(sample_args))

.PHONY: echo-process-args
echo-process-args:
	@echo $(process_args)
	@echo $(process_args) | pbcopy

.PHONY: run-process
run-process:
	@echo "go run main.go process $(process_args)"
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile $$(echo $(process_args))
#	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) go run main.go $$(echo $(process_args))

.PHONY: go-version
go-version:
	@go version

.PHONY: go-env
go-env:
	@go env

.PHONY: release-skip-publish
release-skip-publish: download install-tools generate
	@goreleaser release --skip-publish --snapshot --clean

.PHONY: check-profiler-block
check-profiler-block:
	@go tool pprof -http=:8080 block.pprof

.PHONY: check-profiler-cpu
check-profiler-cpu:
	@go tool pprof -http=:8080 cpu.pprof

.PHONY: check-profiler-memory
check-profiler-memory:
	@go tool pprof -http=:8080 mem.pprof

.PHONY: check-profiler-mutex
check-profiler-mutex:
	@go tool pprof -http=:8080 mutex.pprof

.PHONY: check-profiler-trace
check-profiler-trace:
	@go tool trace -http=:8080 trace.pprof
