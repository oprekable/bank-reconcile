ifeq ($(OS), Windows_NT)
	HELP_CMD = Select-String "^[a-zA-Z_-]+:.*?\#\# .*$$" "./Makefile" | Foreach-Object { $$_data = $$_.matches -split ":.*?\#\# "; $$obj = New-Object PSCustomObject; Add-Member -InputObject $$obj -NotePropertyName ('Command') -NotePropertyValue $$_data[0]; Add-Member -InputObject $$obj -NotePropertyName ('Description') -NotePropertyValue $$_data[1]; $$obj } | Format-Table -HideTableHeaders @{Expression={ $$e = [char]27; "$$e[36m$$($$_.Command)$${e}[0m" }}, Description
else
	HELP_CMD = grep -E '^[a-zA-Z_-]+:.*?\#\# .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?\#\# "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
endif

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@${HELP_CMD}

.PHONY: download
download: ## Download go.mod dependencies
	@echo Download go.mod dependencies
	@go mod download

.PHONY: install-tools
install-tools: download  ## Install required command line tools
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

.PHONY: generate
generate: ## Run go generate google wire dependency injection and mock files
	@export PATH=${PATH}:`go env GOPATH`/bin; wire github.com/oprekable/bank-reconcile/internal/inject
	@export PATH=${PATH}:`go env GOPATH`/bin; go generate ./...

.PHONY: go-lint
go-lint: ## Run golangci-lint (dry run)
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint linters
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint run ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	@export PATH=${PATH}:`go env GOPATH`/bin; staticcheck ./...

.PHONY: govulncheck
govulncheck: ## Run govulncheck to check code vulnerability
	@export PATH=${PATH}:`go env GOPATH`/bin; govulncheck ./...

.PHONY: godeadcode
godeadcode: ## Run deadcode to check dead codes
	@export PATH=${PATH}:`go env GOPATH`/bin; deadcode ./...

.PHONY: development-checks
development-checks: install-tools generate ## Download dependencies, install tools, generate codes, linter, code check (use it in code development cycle)
	@go mod tidy
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint run ./... --fix
	@export PATH=${PATH}:`go env GOPATH`/bin; staticcheck ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; fieldalignment -fix ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; deadcode ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; govulncheck -show verbose ./...

.PHONY: test
test: ## Run unit tests and open coverage page
	@go test -gcflags=all=-l -count=1 -p=8 -parallel=8 -race -coverprofile=coverage.out ./... -json | tee report.json
	@go tool cover -html=coverage.out

.PHONY: run
run: ## Build and run application
	@go build -gcflags -live .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile


UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
	base_args="--from=$$(date -j -v -10d '+%Y-%m-%d') --to=$$(date -j '+%Y-%m-%d')"
endif

ifeq ($(UNAME), Linux)
	base_args="--from=$$(date -d '-10 day' '+%Y-%m-%d') --to=$$(date '+%Y-%m-%d')"
endif

base_args+=" --showlog=true --listbank=bca,bni,mandiri,bri,danamon --profiler=true --debug=true"

process_args="process ${base_args} -s=/tmp/sample/system -b=/tmp/sample/bank -r=/tmp/report"
sample_args="sample ${base_args} --percentagematch=100 --amountdata=100000 -s=/tmp/sample/system -b=/tmp/sample/bank"

.PHONY: echo-sample-args
echo-sample-args: ## Generate command syntax to run application to generate "sample"
	@echo $(sample_args)
	@echo $(sample_args) | pbcopy

.PHONY: run-sample
run-sample: ## Build and run application to generate "sample"
	@echo $(sample_args)
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile  $$(echo $(sample_args))

.PHONY: echo-process-args
echo-process-args: ## Generate command syntax to run application to "process" data
	@echo $(process_args)
	@echo $(process_args) | pbcopy

.PHONY: run-process
run-process: ## Build and run application to "process" data
	@echo "go run main.go process $(process_args)"
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile $$(echo $(process_args))


.PHONY: run-version
run-version: ## Build and run application to show application version
	@echo "go run main.go version"
	@go build -buildvcs=false -ldflags="-s -w" .
	@./bank-reconcile version

.PHONY: go-version
go-version: ## To check current golang version in machine
	@go version

.PHONY: go-env
go-env: ## To check current golang environment variables in machine
	@go env

.PHONY: check-profiler-block
check-profiler-block: ## To open pprof data of block profile
	@go tool pprof -http=:8080 block.pprof

.PHONY: check-profiler-cpu
check-profiler-cpu: ## To open pprof data of cpu profile
	@go tool pprof -http=:8080 cpu.pprof

.PHONY: check-profiler-memory
check-profiler-memory: ## To open pprof data of memory profile
	@go tool pprof -http=:8080 mem.pprof

.PHONY: check-profiler-mutex
check-profiler-mutex: ## To open pprof data of motex profile
	@go tool pprof -http=:8080 mutex.pprof

.PHONY: check-profiler-trace
check-profiler-trace: ## To open pprof data of trace profile
	@go tool trace -http=:8080 trace.pprof
