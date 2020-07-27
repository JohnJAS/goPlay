#!/bin/bash

echo ${PWD}

export GOPATH=/hostgopath

cd ${PWD}/autoUpgrade

GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" ${PWD}/build/output/autoUpgrade