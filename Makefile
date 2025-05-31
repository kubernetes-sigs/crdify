# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL := /usr/bin/env bash -o pipefail
.SHELLFLAGS := -ec

# Tool definitions
GOLANGCI_LINT := go tool -modfile tools/go.mod github.com/golangci/golangci-lint/v2/cmd/golangci-lint

.PHONY: lint
lint: golangci-lint

.PHONY: golangci-lint
golangci-lint:
	$(GOLANGCI_LINT) run $(GOLANGCI_LINT_ARGS)

.PHONY: unit
unit:
	go test ./... -cover -coverprofile unit.out

.PHONY: build
build:
	go build -o bin/crdify main.go

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: verify
verify: fmt tidy lint
	git diff --exit-code

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: e2e
e2e: build
	go run ./test/suite.go

.PHONY: update-e2e
update-e2e: build
	go run ./test/suite.go --update
