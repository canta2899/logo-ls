APP_NAME=logo-ls
BIN_DIR=bin
SRC_DIR=./cmd/logo-ls
BUILD_FLAGS=-ldflags="-s -w" -tags=minimal -trimpath

.PHONY: all bindir clean test install

all: logo-ls

bindir:
	mkdir -p $(BIN_DIR)

logo-ls: bindir
	go build -o $(BIN_DIR)/$(APP_NAME) $(BUILD_FLAGS) $(SRC_DIR)

install:
	go install $(SRC_DIR)

clean:
	rm -rf $(BIN_DIR)

test:
	go test ./...
