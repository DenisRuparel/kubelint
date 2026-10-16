BINARY_NAME := kubelint
BIN_DIR := bin
PKG := ./...

# Detect Go 1.25 toolchain if installed in local go module cache, else system go
GO_TOOLCHAIN := $(wildcard $(HOME)/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.25.0.linux-amd64/bin/go)
ifneq ($(GO_TOOLCHAIN),)
    GO ?= $(GO_TOOLCHAIN)
else
    GO ?= go
endif

.PHONY: all build test vet lint clean fmt help

all: build

help:
	@echo "Available targets:"
	@echo "  build   - Compile kubelint binary into $(BIN_DIR)/$(BINARY_NAME)"
	@echo "  test    - Run unit and integration tests"
	@echo "  vet     - Run go vet static analysis"
	@echo "  lint    - Run golangci-lint (falls back to go vet)"
	@echo "  fmt     - Format Go source code"
	@echo "  clean   - Remove binary and build artifacts"

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build -v -ldflags="-s -w" -o $(BIN_DIR)/$(BINARY_NAME) main.go

test:
	$(GO) test -v $(PKG)

vet:
	$(GO) vet $(PKG)

lint: vet
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed, ran go vet as fallback"; \
	fi

fmt:
	$(GO) fmt $(PKG)

clean:
	rm -rf $(BIN_DIR) dist *.out *.test coverage.html

