# Makefile for building SynapSeq

# Binary information
BIN_NAME 	    := synapseq
BIN_DIR 	    := bin
# Go build metadata
VERSION 	    := $(shell cat VERSION)
HUB_VERSION   	:= $(shell cat HUB_VERSION)
COMMIT  	    := $(shell git rev-parse --short HEAD 2>/dev/null || echo $(shell echo ${GITHUB_SHA} | cut -c1-7))
DATE    	    := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
# Go configuration
GO_METADATA     := -X github.com/ruanklein/synapseq/v3/internal/info.VERSION=$(VERSION) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.BUILD_DATE=$(DATE) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.GIT_COMMIT=$(COMMIT) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.HUB_VERSION=$(HUB_VERSION)
GO_BUILD_FLAGS  := -ldflags="-s -w $(GO_METADATA)"
MAIN 		    := ./cmd/synapseq

.PHONY: all build clean build-windows build-linux build-macos prepare test man install-man

all: build

prepare:
	mkdir -p $(BIN_DIR)

test:
	go test -v ./...

build: prepare
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(MAIN)

build-nohub: prepare
	go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-nohub $(MAIN)

build-windows: prepare
	GOOS=windows GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe $(MAIN)
	GOOS=windows GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64.exe $(MAIN)

build-windows-nohub: prepare
	GOOS=windows GOARCH=amd64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64-nohub.exe $(MAIN)
	GOOS=windows GOARCH=arm64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64-nohub.exe $(MAIN)

build-linux: prepare
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 $(MAIN)
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 $(MAIN)

build-linux-nohub: prepare
	GOOS=linux GOARCH=amd64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64-nohub $(MAIN)
	GOOS=linux GOARCH=arm64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64-nohub $(MAIN)

build-macos: prepare
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64 $(MAIN)

build-macos-nohub: prepare
	GOOS=darwin GOARCH=arm64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64-nohub $(MAIN)

install:
	cp $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

install-nohub:
	cp $(BIN_DIR)/$(BIN_NAME)-nohub /usr/local/bin/$(BIN_NAME)-nohub

clean:
	rm -rf $(BIN_DIR)