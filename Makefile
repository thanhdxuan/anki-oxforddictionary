VERSION ?= $(shell git describe --tags --always --dirty)

TEMPDIR := dist/tmp
MANIFEST := $(TEMPDIR)/manifest.json
REPO_ROOT := $(CURDIR)
ADDON_DIR := $(REPO_ROOT)/addon
BIN_DIR := $(ADDON_DIR)/bin
GOBACKEND_DIR := $(REPO_ROOT)/oxford-lookup
SRCS := $(GOBACKEND_DIR)/cmd/oxford-lookup
GO := go
BUILD_FLAGS := -trimpath -ldflags "-s -w -X main.version=$(VERSION)"
CGO_ENABLED := 0
BINARY_NAME=oxford-lookup
PLATFORMS=linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

ADDON_PACKAGE_NAME := oxford-lookup-$(VERSION).ankiaddon

.PHONY: build test clean pytest be_test package

build:
	@for platform in $(PLATFORMS); do \
		OS=$${platform%/*}; \
		ARCH=$${platform#*/}; \
		OUTPUT=$(BINARY_NAME)-$$OS-$$ARCH; \
		if [ "$$OS" = "windows" ]; then OUTPUT=$$OUTPUT.exe; fi; \
		echo "Building $$OUTPUT..."; \
		CGO_ENABLED=$(CGO_ENABLED) GOOS=$$OS GOARCH=$$ARCH $(GO) build -C $(SRCS) $(BUILD_FLAGS) -o $(BIN_DIR)/$$OUTPUT || exit 1; \
	done

package: clean build $(TEMPDIR)
	cp -r $(BIN_DIR) $(TEMPDIR)
	cp $(ADDON_DIR)/__init__.py $(TEMPDIR)
	cp $(ADDON_DIR)/config.json $(TEMPDIR)
	cp $(ADDON_DIR)/backend.py $(TEMPDIR)
	cp $(ADDON_DIR)/models.py $(TEMPDIR)
	cp $(ADDON_DIR)/rules.py $(TEMPDIR)
	echo "{ \"name\": \"Oxford Audio Lookup\", \"package\": \"oxford_lookup\", \"human_version\": \"$(VERSION)\" }" > $(MANIFEST)
	cd $(TEMPDIR) && zip -r ../$(ADDON_PACKAGE_NAME) .

clean:
	rm -f $(BIN_DIR)/$(BINARY_NAME)-*
	rm -rf $(TEMPDIR) dist/$(ADDON_PACKAGE_NAME)

test: pytest be_test

pytest:
	python -m pytest $(ADDON_DIR)

be_test:
	go test -C $(GOBACKEND_DIR) ./...


$(TEMPDIR):
	mkdir -p $@
