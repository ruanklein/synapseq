# Makefile for building SynapSeq

# Binary information
BIN_NAME 	    := synapseq
BIN_DIR 	    := bin

# Go build metadata
VERSION 	    := $(shell cat VERSION)
HUB_VERSION   	:= $(shell cat HUB_VERSION)
COMMIT  	    := $(shell git rev-parse --short HEAD 2>/dev/null || echo $(shell echo ${GITHUB_SHA} | cut -c1-7))
DATE    	    := $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

# Windows configuration
MAJOR_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f1)
MINOR_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f2)
PATCH_VERSION 			 := $(shell echo $(VERSION) | cut -d. -f3)
GO_VERSION_INFO_CMD 	 := github.com/josephspurrier/goversioninfo/cmd/goversioninfo@v1.5.0
GO_VERSION_INFO_CMD_ARGS := -company="Ruan <ruan.sh>" \
							-description="Synapse-Sequenced Brainwave Generator" \
					  		-copyright="GPL v2" \
					  		-original-name="$(BIN_NAME).exe" \
							-product-name="SynapSeq" \
							-product-version="$(VERSION).0" \
					  		-comment="Main SynapSeq executable" \
							-icon="assets/synapseq.ico" \
							-ver-major=$(MAJOR_VERSION) -product-ver-major=$(MAJOR_VERSION) \
							-ver-minor=$(MINOR_VERSION) -product-ver-minor=$(MINOR_VERSION) \
							-ver-patch=$(PATCH_VERSION) -product-ver-patch=$(PATCH_VERSION) \
							-ver-build=0 -product-ver-build=0
# Go configuration
GO_METADATA     := -X github.com/ruanklein/synapseq/v3/internal/info.VERSION=$(VERSION) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.BUILD_DATE=$(DATE) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.GIT_COMMIT=$(COMMIT) \
				  -X github.com/ruanklein/synapseq/v3/internal/info.HUB_VERSION=$(HUB_VERSION)
GO_BUILD_FLAGS  := -ldflags="-s -w $(GO_METADATA)"
MAIN 		    := ./cmd/synapseq

.PHONY: all build clean test build-windows-amd64 build-windows-arm64 \
		build-windows-nohub-amd64 build-windows-nohub-arm64 \
		build-linux-amd64 build-linux-arm64 \
		build-linux-nohub-amd64 build-linux-nohub-arm64 \
		build-macos build-macos-nohub install install-nohub

all: build

prepare:
	mkdir -p $(BIN_DIR)

# Windows resource file generation
windows-res-amd64:
	go run $(GO_VERSION_INFO_CMD) $(GO_VERSION_INFO_CMD_ARGS) -64 -o cmd/synapseq/synapseq.syso

windows-res-arm64:
	go run $(GO_VERSION_INFO_CMD) $(GO_VERSION_INFO_CMD_ARGS) -arm -o cmd/synapseq/synapseq.syso

test:
	go test -v ./...

build: prepare
	go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME) $(MAIN)

build-nohub: prepare
	go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-nohub $(MAIN)

# Windows builds
build-windows-amd64: prepare windows-res-amd64
	GOOS=windows GOARCH=amd64 go build -tags=windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64.exe $(MAIN)

build-windows-arm64: prepare windows-res-arm64
	GOOS=windows GOARCH=arm64 go build -tags=windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64.exe $(MAIN)

build-windows-nohub-amd64: prepare build-windows-amd64
	GOOS=windows GOARCH=amd64 go build -tags=nohub,windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-amd64-nohub.exe $(MAIN)

build-windows-nohub-arm64: prepare build-windows-arm64
	GOOS=windows GOARCH=arm64 go build -tags=nohub,windows $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-windows-arm64-nohub.exe $(MAIN)


# Linux builds
build-linux-amd64: prepare
	GOOS=linux GOARCH=amd64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64 $(MAIN)

build-linux-arm64: prepare
	GOOS=linux GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64 $(MAIN)

build-linux-nohub-amd64: prepare build-linux-amd64
	GOOS=linux GOARCH=amd64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-amd64-nohub $(MAIN)

build-linux-nohub-arm64: prepare build-linux-arm64
	GOOS=linux GOARCH=arm64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-linux-arm64-nohub $(MAIN)

# macOS builds
build-macos: prepare
	GOOS=darwin GOARCH=arm64 go build $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64 $(MAIN)

build-macos-nohub: prepare
	GOOS=darwin GOARCH=arm64 go build -tags=nohub $(GO_BUILD_FLAGS) -o $(BIN_DIR)/$(BIN_NAME)-macos-arm64-nohub $(MAIN)

# POSIX installation
install:
	cp $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

install-nohub:
	cp $(BIN_DIR)/$(BIN_NAME)-nohub /usr/local/bin/$(BIN_NAME)-nohub

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)
	rm -rf cmd/synapseq/*.syso