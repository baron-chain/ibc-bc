PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION=$(shell git describe --always | sed 's/^v//')
COMMIT=$(shell git log -1 --format='%H')
LEDGER_ENABLED?=true
BUILDDIR?=$(CURDIR)/build
SIMAPP=./testing/simapp
HTTPS_GIT=https://github.com/baron-chain/baron-chain.git

export GO111MODULE=on

# Build tags
build_tags=netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE=$(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags+=ledger
    endif
  else
    UNAME_S=$(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support)
    else
      GCC=$(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags+=ledger
      endif
    endif
  endif
endif

build_tags+=gcc
build_tags+=$(BUILD_TAGS)
build_tags:=$(strip $(build_tags))
build_tags_comma_sep=$(subst $(space),$(comma),$(build_tags))

# Linker flags
ldflags=-X github.com/baron-chain/baron-chain/version.Name=baron \
        -X github.com/baron-chain/baron-chain/version.AppName=barond \
        -X github.com/baron-chain/baron-chain/version.Version=$(VERSION) \
        -X github.com/baron-chain/baron-chain/version.Commit=$(COMMIT) \
        -X "github.com/baron-chain/baron-chain/version.BuildTags=$(build_tags_comma_sep)"

ifeq (rocksdb,$(findstring rocksdb,$(BARON_BUILD_OPTIONS)))
  CGO_ENABLED=1
  BUILD_TAGS+=rocksdb
  ldflags+=-X github.com/baron-chain/baron-chain/types.DBBackend=rocksdb
endif

ifeq (,$(findstring nostrip,$(BARON_BUILD_OPTIONS)))
  ldflags+=-w -s
endif

BUILD_FLAGS=-tags "$(build_tags)" -ldflags '$(ldflags)'
ifeq (,$(findstring nostrip,$(BARON_BUILD_OPTIONS)))
  BUILD_FLAGS+=-trimpath
endif

# Main targets
all: build test

build: BUILD_ARGS=-o $(BUILDDIR)/
build-linux:
	GOOS=linux GOARCH=amd64 LEDGER_ENABLED=false $(MAKE) build

build install: go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

clean:
	rm -rf $(BUILDDIR)/ artifacts/ tmp-swagger-gen/

# Tests
test: test-unit
test-all: test-unit test-ledger-mock test-race test-cover

test-unit: ARGS=-tags='cgo ledger test_ledger_mock'
test-ledger: ARGS=-tags='cgo ledger'
test-race: ARGS=-race -tags='cgo ledger test_ledger_mock'
test-race: TEST_PACKAGES=$(PACKAGES_NOSIMULATION)

run-tests:
	go test -mod=readonly $(ARGS) $(EXTRA_ARGS) $(TEST_PACKAGES)

# Protocol Buffers
protoVer=0.11.6
protoImageName=ghcr.io/baron-chain/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

.PHONY: all build build-linux clean test test-all proto-all proto-gen proto-format proto-lint
