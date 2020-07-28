#!/bin/bash

export https_proxy=web-proxy.us.softwaregrp.net:8080
export http_proxy=web-proxy.us.softwaregrp.net:8080

CURRENT_DIR=${PWD}

WORKSPACE=${CURRENT_DIR}/../../goPlay
OUTPUT_FOLDER=${CURRENT_DIR}/output
GO_VERSION=1.14.6

while [[ ${#} -gt 0 ]]; do
    case "$1" in
        -w|--workspace ) WORKSPACE=${1}
        shift;;
        -o|--output) OUTPUT_FOLDER=${1}
        shift;;
        --go-verison ) GO_VERSION=${1}
        shift;;
    esac
done

#=====================================================================================================#
echo "Pulling golang complie env image..."

docker pull golang:${GO_VERSION}

[[ $? != 0 ]] && "Failed to pull golang build image. Exit" && exit 1

#=====================================================================================================#
echo "Starting to compile CDF upgrade binaries..."

docker run --name build-upgrade-v${GO_VERSION} --privileged --rm -e "OUTPUT_FOLDER=${OUTPUT_FOLDER}" -v ${WORKSPACE}:/workspace -v ${OUTPUT_FOLDER}:/output --workdir /workspace golang:${GO_VERSION} bash build/buildInsideContainer.sh

[[ $? != 0 ]] && "Failed to finish building CDF upgrade binaries" && exit 1