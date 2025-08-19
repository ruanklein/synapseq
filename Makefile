BIN_NAME := synapseq
BIN_DIR := bin
GO_BUILD_FLAGS := -ldflags="-s -w"

.PHONY: all build clean build-windows build-linux build-macos prepare

all: build

prepare:
	mkdir -p $(BIN_DIR)

build: prepare
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME) cmd/main.go

build-windows: prepare
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe cmd/main.go
	GOOS=windows GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64.exe cmd/main.go

build-linux: prepare
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 cmd/main.go
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 cmd/main.go

build-macos: prepare
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64 cmd/main.go

clean:
	rm -rf $(BIN_DIR)