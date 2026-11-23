VERSION ?= $(shell git describe --tags --always --dirty)

REPO_ROOT := $(PWD)
ADDON_DIR := $(REPO_ROOT)/addon
BIN_DIR := $(ADDON_DIR)/bin
GOBACKEND_DIR := $(REPO_ROOT)/oxford-lookup
SRCS := $(GOBACKEND_DIR)/cmd/oxford-lookup
GO := go
BUILD_FLAGS := -trimpath -ldflags "-s -w -X main.version=$(VERSION)"
CGO_ENABLED := 0

BINARY_NAME=oxford-lookup
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

.PHONY: build test clean pytest be_test

build:
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUTPUT=$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building $$OUTPUT..."; \
		CGO_ENABLED=$(CGO_ENABLED) GOOS=$$OS GOARCH=$$ARCH $(GO) build -C $(SRCS) $(BUILD_FLAGS) -o $(BIN_DIR)/$$OUTPUT || exit 1; \
	done

clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)-*

test: pytest be_test

pytest:
	python -m pytest $(ADDON_DIR)

be_test:
	go test -C $(GOBACKEND_DIR) ./...
	
