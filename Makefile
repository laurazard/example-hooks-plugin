ifeq ($(OS),Windows_NT)
    DETECTED_OS = Windows
else
    DETECTED_OS = $(shell uname -s)
endif

ifeq ($(DETECTED_OS),Windows)
	BINARY_EXT=.exe
endif

DESTDIR ?= ./bin/build


all: build

.PHONY: build ## Build the hints cli-plugin
build:
	GO111MODULE=on go build $(BUILD_FLAGS) -trimpath -tags "$(GO_BUILDTAGS)" -ldflags "$(GO_LDFLAGS)" -o "$(DESTDIR)/docker-hints$(BINARY_EXT)" ./cmd

.PHONY: install
install: build
	mkdir -p ~/.docker/cli-plugins
	install bin/build/docker-hints ~/.docker/cli-plugins/docker-hints
