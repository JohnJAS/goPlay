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

check:
	@echo "UNAME_S: ${UNAME_S}, OS: ${OS}, ARCH: ${ARCH}"

compile:
	@echo "Compiling linux version..."
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" .
	@echo "Compiling windows version..."
	GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" .
