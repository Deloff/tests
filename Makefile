ifneq ($(wildcard .env),)
	include .env
endif

GOOS=linux
GOARCH=amd64
GOPRIVATE=github.com
CGO_ENABLED=0

PROJECT_DIR = $(shell pwd)
PROJECT_BIN = $(PROJECT_DIR)/bin
PROJECT_TMP = $(PROJECT_DIR)/tmp

$(shell [ -f $(PROJECT_BIN) ] || mkdir -p $(PROJECT_BIN))


gci:
	go install github.com/daixiang0/gci@latest

gofumpt:
	go install mvdan.cc/gofumpt@latest

fmt: gci gofumpt
	gci write -s standard -s default -s "prefix(gitlab.mk-dev.ru)" . --skip-generated
	gofumpt -e -w -extra .

.dev-tools:
	@[ -f $(GOLANGCI_LINT_BIN) ] || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(PROJECT_BIN) v1.55.1

GOLANGCI_LINT_BIN = $(PROJECT_BIN)/golangci-lint
GOLANGCI_LINT_CONFIG = $(PROJECT_DIR)/.golangci.yaml
lint: .dev-tools
	$(GOLANGCI_LINT_BIN) run $(PROJECT_DIR)/... --config=$(GOLANGCI_LINT_CONFIG)

