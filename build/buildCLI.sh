#!/bin/bash

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

CURRENT_DIR=${PWD}
WORKSPACE=${CURRENT_DIR}/../
SSH_CONFIG_FOLDER=${HOME}/.ssh/
GO_VERSION=1.14.6

while [[ ${#} -gt 0 ]]; do
    case "$1" in
        -w|--workspace ) WORKSPACE=${1}
        shift;;
        --ssh-config ) SSH_CONFIG_FOLDER=${1}
        shift;;
        --go-verison ) GO_VERSION=${1}
        shift;;
    esac
done

#=====================================================================================================#
echo "Pulling golang complie env image..."

docker pull golang:${GO_VERSION}

[[ $? != 0 ]] && "Failed to pull golang build image. Exit" && exit 1

#=============================================================root========================================#
echo "Starting to compile CDF upgrade binaries..."

#debug : docker run --name build-upgrade-v${GO_VERSION} --privileged --rm -v ${WORKSPACE}:/workspace -v ${OUTPUT_FOLDER}:/output -v ${SSH_CONFIG_FOLDER}:/root/.ssh --workdir /workspace -ti golang:${GO_VERSION} bash 
docker run --name build-upgrade-v${GO_VERSION} --privileged --rm -v ${WORKSPACE}:/workspace -v ${OUTPUT_FOLDER}:/output -v ${SSH_CONFIG_FOLDER}:/root/.ssh --workdir /workspace golang:${GO_VERSION} bash build/buildInsideContainer.sh

[[ $? != 0 ]] && "Failed to finish building CDF upgrade binaries" && exit 1