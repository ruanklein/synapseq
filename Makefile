BIN_NAME := synapseq
BIN_DIR := bin
VERSION := $(shell cat VERSION)
GO_BUILD_FLAGS := -ldflags="-s -w -X github.com/ruanklein/synapseq/internal/info.VERSION=$(VERSION)"
MAIN := cmd/synapseq/main.go
# Documentation
MAN_DIR := docs/manpage
MAN_FILE := $(MAN_DIR)/synapseq.1
MAN_INSTALL_DIR := /usr/local/share/man/man1

.PHONY: all build clean build-windows build-linux build-macos prepare test man install-man

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

install:
	cp $(BIN_DIR)/$(BIN_NAME) /usr/local/bin/$(BIN_NAME)

# Documentation
man:
	@echo "Generating man page..."
	@mkdir -p $(MAN_DIR)
	@pandoc docs/USAGE.md -s -t man \
		-V title="SYNAPSEQ" \
		-V section="1" \
		-V header="SynapSeq Manual" \
		-V footer="SynapSeq 3.1" \
		-V date="$$(date +'%B %Y')" \
		-o $(MAN_FILE)
	@echo "Man page generated at $(MAN_FILE)"

install-man: man
	@echo "Installing man page..."
	@mkdir -p $(MAN_INSTALL_DIR)
	@cp $(MAN_FILE) $(MAN_INSTALL_DIR)/$(BIN_NAME).1
	@gzip -f $(MAN_INSTALL_DIR)/$(BIN_NAME).1
	@echo "Man page installed to $(MAN_INSTALL_DIR)/$(BIN_NAME).1.gz"
	@echo "You can now use: man $(BIN_NAME)"

clean:
	rm -rf $(BIN_DIR)
	rm -rf $(MAN_DIR)