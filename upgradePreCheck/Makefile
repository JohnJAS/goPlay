.PHONY: all compile

all: compile

compile:
	@echo "Compiling linux version..."
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" .
