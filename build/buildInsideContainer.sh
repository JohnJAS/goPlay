#!/bin/bash

echo ${PWD}

cd ${PWD}/autoUpgrade

GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" ${PWD}/build/output/autoUpgrade