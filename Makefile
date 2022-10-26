.DEFAULT_GOAL 	:= help
.PHONY: compile run run-cli test coverage clients logs all gomod_tidy go_fmt help

# Variables
SEDGE_VERSION = $(shell git tag | sort | tail -n 1)
LDFLAGS=-X github.com/NethermindEth/sedge/internal/utils.Version="${SEDGE_VERSION}"

# Commands
compile: ## compile:
	@mkdir -p build
	@go build -ldflags "${LDFLAGS}" -o build/sedge cmd/sedge/main.go

compile-linux: ## compile:
	@mkdir -p build
	@env GOOS=linux go build -ldflags="${LDFLAGS[*]}" -o build/sedge cmd/main.go

install: compile ## compile the binary and copy it to PATH
	@sudo cp build/sedge /usr/local/bin

run: ## run
	@./build/sedge

run-cli: compile ## run cli
	@./build/sedge cli --config ./config.yaml

test: ## run tests
	@mkdir -p coverage
	@go test -coverprofile=coverage/coverage.out -covermode=count ./...

codecov-test: ## unit tests with coverage using the courtney tool
	@mkdir -p coverage
	@courtney/courtney -v -o coverage/coverage.out ./...

install-gofumpt: ## install gofumpt
	go install mvdan.cc/gofumpt@latest

install-courtney: ## Install courtney for code coverage
	@git clone https://github.com/stdevMac/courtney
	@(cd courtney && go get  ./... && go build courtney.go)
	@go get ./...

install-deps: | install-gofumpt install-courtney ## Install some project dependencies

coverage: coverage/coverage.out ## show tests coverage
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html

clients: compile ## Run cmd clients with ./config.yaml
	@./build/sedge clients --config ./config.yaml

logs: compile ## run cmd logs with ./config.yaml
	@./build/sedge logs --config ./config.yaml

all: compile run ## build and run

gomod_tidy: ## go mod tidy
	 go mod tidy

format: ## run code formatting
	gofumpt -l -w .

format-check: ## check formatting
	# assert `gofumpt -l` produces no output
	test ! $$(gofumpt -l . | tee /dev/stderr)


help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'