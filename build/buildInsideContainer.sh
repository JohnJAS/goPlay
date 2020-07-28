#!/bin/bash

CURRENT_DIR=${PWD}
[[ -z $OUTPUT_FOLDER ]] && OUTPUT_FOLDER=${CURRENT_DIR}/build/output

git config --global url."git@github.com:".insteadOf "https://github.com/"
go env -w GOPROXY=https://goproxy.cn/
go env -w GOPRIVATE=github.com/JohnJAS/*

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

cd ${CURRENT_DIR}/autoUpgrade

echo "Building autoUpgrade linux version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o $OUTPUT_FOLDER/autoUpgrade/autoUpgrade

[[ $? == 0 ]] && echo "success" || echo "failed"

echo "Building autoUpgrade windows version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o $OUTPUT_FOLDER/autoUpgrade/autoUpgrade.exe
[[ $? == 0 ]] && echo "success" || echo "failed"

echo "Building upgradePreCheck"
GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o $OUTPUT_FOLDER/upgradePreCheck/upgradePreCheck
[[ $? == 0 ]] && echo "success" || echo "failed"