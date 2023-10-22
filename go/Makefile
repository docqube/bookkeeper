######################################################
# config
######################################################
# setting SHELL to bash allows bash commands to be executed by recipes
# options are set to exit when a recipe line exits non-zero or a piped command fails
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

## location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin

# golang ci linter binary and version
GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.51.2

# tparse binary and version
TPARSE ?= $(LOCALBIN)/tparse
TPARSE_VERSION ?= latest


######################################################
# misc
######################################################
.PHONY: clean
clean:
	rm -rf build

.PHONY: localbin
localbin:
	mkdir -p $(LOCALBIN)

.PHONY: golangci-lint
golangci-lint: localbin
	test -s $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) || GOBIN=$(LOCALBIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	test -s $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) || mv $(GOLANGCI_LINT) $(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION)

.PHONY: tparse
tparse: localbin
	test -s $(TPARSE) || GOBIN=$(LOCALBIN) go install github.com/mfridman/tparse@$(TPARSE_VERSION)
	test -s $(TPARSE)-$(TPARSE_VERSION) || mv $(TPARSE) $(TPARSE)-$(TPARSE_VERSION)


######################################################
# go
######################################################
.PHONY: tidy
tidy:
	go mod tidy

.PHONY: vendor
vendor:
	go mod vendor

.PHONY: build
build:
	go build -v -tags netgo -o build/bookkeeper-api cmd/api/main.go

.PHONY: run
run:
	go run ./cmd/api/main.go

.PHONY: update
update:
	go get -u ./...
	$(MAKE) tidy
	$(MAKE) vendor


######################################################
# lint
######################################################
.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT)-$(GOLANGCI_LINT_VERSION) run ./... --timeout 5m0s -v


######################################################
# test
######################################################
.PHONY: test
test: tparse
	set -eu
	set -o pipefail
	go test ./... -cover -json | $(TPARSE)-$(TPARSE_VERSION) -all
