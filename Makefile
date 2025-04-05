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
	@cat tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %@latest

.PHONY: generate
generate:
	@export PATH=${PATH}:`go env GOPATH`/bin; wire github.com/oprekable/bank-reconcile/internal/inject
	@go generate ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; mockery

.PHONY: go-lint
go-lint:
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint linters
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint run ./...

.PHONY: staticcheck
staticcheck:
	@export PATH=${PATH}:`go env GOPATH`/bin; staticcheck ./...

.PHONY: govulncheck
govulncheck:
	@export PATH=${PATH}:`go env GOPATH`/bin; govulncheck ./...

.PHONY: godeadcode
godeadcode:
	@export PATH=${PATH}:`go env GOPATH`/bin; deadcode ./...

.PHONY: go-lint-fix-struct-staticcheck-govulncheck
go-lint-fix-struct-staticcheck-govulncheck: install-tools generate
	@go mod tidy
	@export PATH=${PATH}:`go env GOPATH`/bin; golangci-lint run ./... --fix
	@export PATH=${PATH}:`go env GOPATH`/bin; staticcheck ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; fieldalignment -fix ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; deadcode ./...
	@export PATH=${PATH}:`go env GOPATH`/bin; govulncheck -show verbose ./...

.PHONY: test
test:
	@go test -gcflags=all=-l -count=1 -p=8 -parallel=8 -race -coverprofile=coverage.out ./... -json | tee report.json
	@go tool cover -html=coverage.out

.PHONY: run
run:
	@go build -gcflags -live .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile



UNAME := $(shell uname)
ifeq ($(UNAME), Darwin)
	base_args="--showlog=true --listbank=bca,bni,mandiri,bri,danamon --from=$$(date -j -v -10d '+%Y-%m-%d') --to=$$(date -j '+%Y-%m-%d')"
endif

ifeq ($(UNAME), Linux)
	base_args="--showlog=true --listbank=bca,bni,mandiri,bri,danamon --from=$$(date -d '-10 day' '+%Y-%m-%d') --to=$$(date '+%Y-%m-%d')"
endif

base_args+=" -i=true -g=false"

process_args="process ${base_args} -s=/tmp/sample/system -b=/tmp/sample/bank -r=/tmp/report"
sample_args="sample ${base_args} --percentagematch=100 --amountdata=100000 -s=/tmp/sample/system -b=/tmp/sample/bank"

.PHONY: echo-sample-args
echo-sample-args:
	@echo $(sample_args)
	@echo $(sample_args) | pbcopy

.PHONY: run-sample
run-sample:
	@echo $(sample_args)
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile  $$(echo $(sample_args))

.PHONY: echo-process-args
echo-process-args:
	@echo $(process_args)
	@echo $(process_args) | pbcopy

.PHONY: run-process
run-process:
	@echo "go run main.go process $(process_args)"
	@go build -buildvcs=false -ldflags="-s -w" .
	@env $$(cat "params/.env" | grep -Ev '^#' | xargs) ./bank-reconcile $$(echo $(process_args))


.PHONY: run-version
run-version:
	@echo "go run main.go version"
	@go build -buildvcs=false -ldflags="-s -w" .
	@./bank-reconcile version

.PHONY: go-version
go-version:
	@go version

.PHONY: go-env
go-env:
	@go env

.PHONY: release-skip-publish
release-skip-publish: download install-tools generate
	@export PATH=${PATH}:`go env GOPATH`/bin; goreleaser release --skip-publish --snapshot --clean

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
