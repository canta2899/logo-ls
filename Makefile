APP_NAME=logo-ls
SRC_DIR=./cmd/logo-ls
VERSION=$$(git describe --tags --abbrev=0 --always)
BUILD_FLAGS=-ldflags="-s -w -X 'github.com/canta2899/logo-ls/app.Version=$(VERSION)'" -tags=minimal -trimpath

OUTPUT_NAME=$(APP_NAME)$(if $(findstring windows,$(GOOS)),.exe,)

GORELEASER_VERSION ?= latest
GORELEASER=go run github.com/goreleaser/goreleaser/v2@$(GORELEASER_VERSION)

.PHONY: all bindir clean test test-clean install logo-ls release-check release-snapshot

all: logo-ls

logo-ls:
	go build -o $(OUTPUT_NAME) $(BUILD_FLAGS) $(SRC_DIR)


install:
	go install $(SRC_DIR)

clean:
	rm logo-ls

test:
	go test ./...

test-clean:
	go clean -testcache

release-check:
	$(GORELEASER) check

release-snapshot:
	$(GORELEASER) release --snapshot --clean
