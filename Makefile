BIN_NAME := synapseq
BIN_DIR := bin
GO_BUILD_FLAGS := -ldflags="-s -w"
MAIN := cmd/synapseq/main.go

.PHONY: all build clean build-windows build-linux build-macos prepare test

all: build

prepare:
	mkdir -p $(BIN_DIR)

test:
	go test -v ./...

build: prepare
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(MAIN)

build-windows: prepare
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe $(MAIN)
	GOOS=windows GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64.exe $(MAIN)

build-linux: prepare
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 $(MAIN)
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 $(MAIN)

build-macos: prepare
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64 $(MAIN)

clean:
	rm -rf $(BIN_DIR)