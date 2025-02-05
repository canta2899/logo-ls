APP_NAME=logo-ls
SRC_DIR=./cmd/logo-ls
VERSION=$$(git describe --tags --abbrev=0)
BUILD_FLAGS=-ldflags="-s -w -X 'github.com/canta2899/logo-ls/app.Version=$(VERSION)'" -tags=minimal -trimpath

OUTPUT_NAME=$(APP_NAME)$(if $(findstring windows,$(GOOS)),.exe,)

.PHONY: all bindir clean test install

all: logo-ls

logo-ls:
	go build -o $(OUTPUT_NAME) $(BUILD_FLAGS) $(SRC_DIR)


install:
	go install $(SRC_DIR)

clean:
	rm logo-ls

test:
	go test ./...
