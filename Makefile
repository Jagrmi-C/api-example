default: help

.PHONY: help
help: # Show help for each of the Makefile recipes.
	@grep -E '^[a-zA-Z0-9 -]+:.*#'  Makefile | sort | while read -r l; do printf "\033[1;32m$$(echo $$l | cut -f 1 -d':')\033[00m:$$(echo $$l | cut -f 2- -d'#')\n"; done

VERSION_TAG := $(shell git describe --tags --always)
VERSION_VERSION := $(shell git log --date=iso --pretty=format:"%cd" -1) $(VERSION_TAG)
VERSION_COMPILE := $(shell date +"%F %T %z") by $(shell go version)
VERSION_BRANCH  := $(shell git rev-parse --abbrev-ref HEAD)
VERSION_GIT_DIRTY := $(shell git diff --no-ext-diff 2>/dev/null | wc -l | awk '{print $1}')
VERSION_DEV_PATH:= $(shell pwd)

# Go Checkup
GOPATH ?= $(shell go env GOPATH)
GO111MODULE:=auto
export GO111MODULE
ifeq "$(GOPATH)" ""
	$(error Please set the environment variable GOPATH before running `make`)
endif

PATH := ${GOPATH}/bin:$(PATH)
GCFLAGS=-gcflags="all=-trimpath=${GOPATH}"
LDFLAGS=-ldflags="-s -w -X 'main.Version=${VERSION_VERSION}' -X 'main.Compile=${VERSION_COMPILE}' -X 'main.Branch=${VERSION_BRANCH}'"

GO = go

FILES_TO_FMT      ?= $(shell find . -path ./vendor -prune -o \( -name "*.go" ! -name "*_moq.go" \) -print)

# Commands
.PHONY: all
all: | init deps

.PHONY: init
init: ; $(info $(M) Installing tools dependencies ...) @ ## Install tools dependencies
	pip3 install -U commitizen pre-commit

	pre-commit install --install-hooks
	pre-commit install --hook-type commit-msg

.PHONY: deps
deps: ## Ensures fresh go.mod and go.sum.
	@echo ">> running checking deps commands"
	@go mod tidy
	@go mod verify

.PHONY: lint
lint:  ; $(info $(M) Running linting process ...)	@ ## Run golangci linter
	$Q golangci-lint run --timeout 5m0s -v --out-format colored-line-number ./...
	@echo "🚀 Lint stage Completed"

.PHONY: go-format
go-format: ## Formats Go code including imports.
go-format: $(GOIMPORTS) deps
	@echo ">> formatting go code"
	@gofmt -s -w $(FILES_TO_FMT)
	@$(GOIMPORTS) -w $(FILES_TO_FMT)

.PHONY: mocks
mocks: ## Generates mocks for all
	mockery
	@echo ">> Generating mocks completed"

.PHONY: vuln
vuln: # Vulnerability Management for Go
	$Q govulncheck ./...
	@echo "🧪 Vulnerability Checkup Completed"

.PHONY: test
test: ## Run all tests
	@go test -v -coverprofile coverage.out ./...
	@echo "🧪 Testing Stage Completed"

.PHONY: test-unit
test-unit: ## Run unit-tests
	@go test -v -race -coverprofile coverage.out -short ./...
	@echo "🧪 Unit-Testing Stage Completed"

.PHONY: validate
validate: lint test vuln ## Run complex process to verify ability to push changes
	@echo "🧪 Validation the ability to Push Completed"

.PHONY: fieldalignment-fix
fieldalignment-fix: ; $(info $(M) Running the linter fieldalignment to fix issues ...) @ ## Fix fieldalignment e.g. ./auth0/auth
	$Q fieldalignment -fix -test=false $(ARG)
	@echo "📦 Fix structure aligment completed"

.PHONY: build
build: ; $(info $(M) Building executable...) @ ## Build program binary
	$Q mkdir -p bin
	$Q ret=0 && for d in $$($(GO) list -f '{{if (eq .Name "main")}}{{.ImportPath}}{{end}}' ./...); do \
		b=$$(basename $${d}) ; \
		$(GO) build ${LDFLAGS} ${GCFLAGS} -o bin/$${b} $$d || ret=$$? ; \
		echo "$(M) Build: bin/$${b}" ; \
		echo "$(M) Done!" ; \
	done ; exit $$ret

tools-lint:
	$Q $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2

tools:
	$Q $(GO) install github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@latest
	$Q $(GO) install github.com/vektra/mockery/v2@v2.40.1
