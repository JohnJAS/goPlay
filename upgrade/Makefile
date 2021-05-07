BINDIR      := $(CURDIR)/bin
INSTALL_PATH ?= /usr/local/bin
BINNAME     ?= upgrade

# go option
PKG        := ./...
TAGS       :=
LDFLAGS    := -w
GOFLAGS    :=

# Required for globs to work correctly
SHELL      = /usr/bin/env bash

GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_SHA    = $(shell git rev-parse --short HEAD)
GIT_TAG    = $(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY  = $(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

ifdef VERSION
	BINARY_VERSION = $(VERSION)
endif
BINARY_VERSION ?= ${GIT_TAG}

# Only set Version if building a tag or VERSION is set
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X upgrade/internal/version.version=${BINARY_VERSION}
endif

VERSION_METADATA = unreleased
# Clear the "unreleased" string in BuildMetadata
ifneq ($(GIT_TAG),)
	VERSION_METADATA =
endif

LDFLAGS += -X upgrade/internal/version.metadata=${VERSION_METADATA}
LDFLAGS += -X upgrade/internal/version.gitCommit=${GIT_COMMIT}
LDFLAGS += -X upgrade/internal/version.gitTreeState=${GIT_DIRTY}

.PHONY: all
all: build

# ------------------------------------------------------------------------------
#  build
.PHONY: build
build: $(BINDIR)/$(BINNAME)

$(BINDIR)/$(BINNAME):
	GO111MODULE=on go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o '$(BINDIR)'/$(BINNAME) ./cmd

.PHONY: clean
clean:
	@rm -rf '$(BINDIR)' ./_dist

.PHONY: info
info:
	 @echo "Version:           ${VERSION}"
	 @echo "Git Tag:           ${GIT_TAG}"
	 @echo "Git Commit:        ${GIT_COMMIT}"
	 @echo "Git Tree State:    ${GIT_DIRTY}"
