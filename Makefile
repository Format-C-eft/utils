export GO111MODULE=on

ENV_PARAMS ?= CGO_ENABLED=0
LOCAL_BIN := $(CURDIR)/bin
LINT_BIN := $(LOCAL_BIN)/golangci-lint

.PHONY: deps
deps:
	$(info Download dependencies...)
	go mod download

.PHONY: test
test:
	$(info Running test...)
	go test ./...

.PHONY: .lint
.lint:
	$(info Running lint...)
	$(LINT_BIN) run --new-from-rev=origin/master --config=.golangci.pipeline.yaml ./...

.PHONY: .lint-full
.lint-full:
	$(info Running lint-full...)
	$(LINT_BIN) run --config=.golangci.pipeline.yaml ./...

.PHONY: lint
lint: install-lint .lint

.PHONY: lint-full
lint-full: install-lint .lint-full

LINT_TAG ?= 1.55.2
install-lint: export GOBIN := $(LOCAL_BIN)
install-lint:
	$(info Installing golangci-lint v$(LINT_TAG))
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(LINT_TAG)
	go mod tidy