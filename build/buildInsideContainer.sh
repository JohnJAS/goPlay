#!/bin/bash

echo ${PWD}

go env -w GOPROXY=https://goproxy.cn/

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

cd ${PWD}/autoUpgrade

GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w" ${PWD}/build/output/autoUpgrade