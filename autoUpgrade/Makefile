.PHONY: all compile

ifeq ($(OS),Windows_NT)
  OS := windows
  ARCH := amd64
else
  UNAME_S := $(shell uname -s)
  ifeq ($(UNAME_S),Linux)
    OS := linux
    ARCH := amd64
  endif
endif

all: compile

debug:
	@echo "UNAME_S: ${UNAME_S}, OS: ${OS}, ARCH: ${ARCH}"

compile:
	@echo "Compiling on ${OS}..."
	GO111MODULE=on CGO_ENABLED=0 GOOS=${OS} GOARCH=${ARCH} go build -ldflags "-w" .
