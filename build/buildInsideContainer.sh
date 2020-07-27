#!/bin/bash

CURRENT_DIR=${PWD}

go env -w GOPROXY=https://goproxy.cn/

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

cd ${CURRENT_DIR}/autoUpgrade

echo "Building autoUpgrade linux version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o ${CURRENT_DIR}/build/output/autoUpgrade

[[ $? == 0 ]] && echo "success" || echo "failed"

echo "Building autoUpgrade windows version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o ${CURRENT_DIR}/build/output/autoUpgrade.exe
[[ $? == 0 ]] && echo "success" || echo "failed"