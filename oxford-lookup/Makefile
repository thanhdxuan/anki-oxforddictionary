GOOS ?=
GOARCH ?=

TARGET := app
SRCS := ./cmd/oxford-lookup/
GO := go
GOFLAGS := -ldflags="-s -w"
BIN := oxford-lookup-$(GOOS)-$(GOARCH)
BUILD_CMD := $(GO) build

$(TARGET):
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(BUILD_CMD) -o $(BIN) $(SRCS)
