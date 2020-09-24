PACKAGES_NOSIMULATION=$(shell go list ./... | grep -v '/simulation')
BINDIR ?= $(GOPATH)/bin
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')

export GO111MODULE = on

LEDGER_ENABLED ?= true

########################################
### Build tags

build_tags = netgo

ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=lendingD \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=lendingCLI \
		  -X github.com/cosmos/cosmos-sdk/version.Name=lendingnetwork \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

all: build lint test

########################################
### Install

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/lendingD
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/lendingCLI

########################################
### Build

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly -o ./build/lendingD.exe $(BUILD_FLAGS) ./cmd/lendingD
	go build -mod=readonly -o ./build/lendingCLI.exe $(BUILD_FLAGS) ./cmd/lendingCLI
else
	go build -mod=readonly -o ./build/lendingD $(BUILD_FLAGS) ./cmd/lendingD
	go build -mod=readonly -o ./build/lendingCLI $(BUILD_FLAGS) ./cmd/lendingCLI
endif


build-darwin: go.sum
	env GOOS=darwin GOARCH=amd64 go build -mod=readonly -o ./build/Darwin-AMD64/lendingCLI $(BUILD_FLAGS) ./cmd/lendingCLI
	env GOOS=darwin GOARCH=amd64 go build -mod=readonly -o ./build/Darwin-AMD64/lendingD $(BUILD_FLAGS) ./cmd/lendingD

build-linux: go.sum
	env GOOS=linux GOARCH=amd64 go build -mod=readonly -o ./build/Linux-AMD64/lendingCLI $(BUILD_FLAGS) ./cmd/lendingCLI
	env GOOS=linux GOARCH=amd64 go build -mod=readonly -o ./build/Linux-AMD64/lendingD $(BUILD_FLAGS) ./cmd/lendingD

build-windows: go.sum
	env GOOS=windows GOARCH=amd64 go build -mod=readonly -o ./build/Windows-AMD64/lendingCLI.exe $(BUILD_FLAGS) ./cmd/lendingCLI
	env GOOS=windows GOARCH=amd64 go build -mod=readonly -o ./build/Windows-AMD64/lendingD.exe $(BUILD_FLAGS) ./cmd/lendingD

build-all: go.sum
	make build-darwin
	make build-linux
	make build-windows

prepare-release: go.sum build-all
	rm -f ./build/Darwin-386.zip ./build/Darwin-AMD64.zip
	rm -f ./build/Linux-386.zip ./build/Linux-AMD64.zip
	rm -f ./build/Windows-386.zip ./build/Windows-AMD64.zip
	zip -jr ./build/Darwin-AMD64.zip ./build/Darwin-AMD64/lendingCLI ./build/Darwin-AMD64/lendingD
	zip -jr ./build/Linux-AMD64.zip ./build/Linux-AMD64/lendingCLI ./build/Linux-AMD64/lendingD
	zip -jr ./build/Windows-AMD64.zip ./build/Windows-AMD64/lendingCLI.exe ./build/Windows-AMD64/lendingD.exe

########################################
### Tools & dependencies

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	go mod verify

lint:
	golangci-lint run
	find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

########################################
### Testing

test: test_unit

test_unit:
	@VERSION=$(VERSION) go test -mod=readonly $(PACKAGES_NOSIMULATION) -tags='ledger test_ledger_mock'

.PHONY: lint test test_unit go-mod-cache

clean:
	rm -rf build/

.PHONY: clean
