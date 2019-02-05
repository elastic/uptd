export VERSION := v1.0.0

GOLINT_PRESENT := $(shell command -v golint 2> /dev/null)
GOIMPORTS_PRESENT := $(shell command -v goimports 2> /dev/null)
GOLICENSER_PRESENT := $(shell command -v go-licenser 2> /dev/null)
TEST_UNIT_FLAGS ?= -timeout 10s -p 4 -race -cover
TEST_UNIT_PACKAGE ?= ./...

.PHONY: deps
deps:
ifndef GOLINT_PRESENT
	@ go get -u golang.org/x/lint/golint
endif
ifndef GOIMPORTS_PRESENT
	@ go get -u golang.org/x/tools/cmd/goimports
endif
ifndef GOLICENSER_PRESENT
	@ go get -u github.com/elastic/go-licenser
endif

.PHONY: lint
lint:
	@ golint -set_exit_status $(shell go list ./...)
	@ gofmt -d -e -s .
	@ go-licenser -d

.PHONY: format
format: deps
	@ gofmt -e -w -s .
	@ goimports -w .
	@ go-licenser

.PHONY: unit
unit:
	@ go test $(TEST_UNIT_FLAGS) $(TEST_UNIT_PACKAGE)

.PHONY: tag
tag:
	@ git tag $(VERSION)

.PHONY: release
release:
	@ git push upstream $(VERSION)
