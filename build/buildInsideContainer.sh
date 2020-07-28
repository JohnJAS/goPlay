#!/bin/bash

CURRENT_DIR=${PWD}

git config --global url."git@github.com:".insteadOf "https://github.com/"
go env -w GOPROXY=https://goproxy.cn/
go env -w GOPRIVATE=github.com/JohnJAS/*

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

#===============================================================================================================================#
cd ${CURRENT_DIR}/autoUpgrade

echo "Building autoUpgrade linux version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" -o ${CURRENT_DIR}/build/output/autoUpgrade/autoUpgrade

[[ $? == 0 ]] && echo "success" || echo "failed"

echo "Building autoUpgrade windows version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o ${CURRENT_DIR}/build/output/autoUpgrade/autoUpgrade.exe
[[ $? == 0 ]] && echo "success" || echo "failed"

#===============================================================================================================================#
cd ${CURRENT_DIR}/upgradePreCheck

echo "Building upgradePreCheck only linux version..."
GO111MODULE=on CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-w" -o ${CURRENT_DIR}/build/output/upgradePreCheck/upgradePreCheck
[[ $? == 0 ]] && echo "success" || echo "failed"