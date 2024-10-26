SHELL := /bin/bash

# Environment detection
GO        := $(shell which go 2>/dev/null)
UNAME_S   := $(shell uname -s 2>/dev/null)
UNAME_M   := $(shell uname -m 2>/dev/null)
GOPATH    := $(shell go env GOPATH)

# OS-specific settings
ifeq ($(OS),Windows_NT)
    FS      := \\
    GO      := $(shell where go.exe 2>NUL)
    NULL    := NUL
    MKDIR   := mkdir
    RMDIR   := rmdir /s /q
else
    FS      := /
    NULL    := /dev/null
    MKDIR   := mkdir -p
    RMDIR   := rm -rf
endif

# Directories and paths
PREFIX       ?= /usr/local
BIN          ?= $(PREFIX)/bin
TOOLS_DIR    := $(GOPATH)/bin
GITHUB_DIR   := $(GOPATH)$(FS)src$(FS)github.com
MKFILE_PATH  := $(abspath $(lastword $(MAKEFILE_LIST)))
MKFILE_DIR   := $(dir $(MKFILE_PATH))

# Tool versions and paths
BUF_VERSION  ?= 0.11.0
STATIK       := $(TOOLS_DIR)/statik
RUNSIM       := $(TOOLS_DIR)/runsim

# Validate Go installation
ifeq ($(GO),)
    $(error Go not found. Please install Go and ensure it's in PATH)
endif

# Function to safely create GitHub directory and clone/update repository
define fetch_repo
	@echo "Fetching $(1)/$(2)@$(3)..."
	@$(MKDIR) "$(GITHUB_DIR)$(FS)$(1)" 2>$(NULL) || true
	@if [ ! -d "$(GITHUB_DIR)$(FS)$(1)$(FS)$(2)" ]; then \
		cd "$(GITHUB_DIR)$(FS)$(1)" && git clone https://github.com/$(1)/$(2); \
	fi
	@cd "$(GITHUB_DIR)$(FS)$(1)$(FS)$(2)" && git fetch origin && git checkout -q $(3)
endef

.PHONY: all tools clean statik runsim

all: tools

# Tools installation
tools: tools-stamp

tools-stamp: statik runsim
	@touch $@

statik: $(STATIK)
$(STATIK):
	@echo "Installing statik..."
	@go install github.com/rakyll/statik@v0.1.6
	@which statik > /dev/null || { echo "statik installation failed"; exit 1; }

runsim: $(RUNSIM)
$(RUNSIM):
	@echo "Installing runsim..."
	@go install github.com/cosmos/tools/cmd/runsim@v1.0.0
	@which runsim > /dev/null || { echo "runsim installation failed"; exit 1; }

# Cleaning targets
clean: tools-clean

tools-clean:
	@echo "Cleaning tools..."
	@rm -f $(STATIK) $(RUNSIM)
	@rm -f tools-stamp

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all          - Install all tools"
	@echo "  tools        - Install development tools"
	@echo "  statik       - Install statik tool"
	@echo "  runsim       - Install runsim tool"
	@echo "  clean        - Clean all generated files"
	@echo "  tools-clean  - Clean installed tools"
	@echo "  help         - Show this help message"
