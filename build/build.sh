#!/bin/bash

CURRENT_DIR=${PWD}

WORKSPACE=${CURRENT_DIR}/../../goPlay
OUTPUT_FOLDER=${CURRENT_DIR}/output
OS=linux
ARCH=amd64
GO_VERSION=1.14.6

while [[ ${#} -gt 0 ]]; do 
    case "$1" in
        -w|--workspace ) WORKSPACE=${1}
        shift;;
        -o|--output) OUTPUT_FOLDER=${1}
        shift;;
        --os ) OS=${1}
        shift;;
        --arch ) ARCH=${1}
        shift;;
        --go-verison ) GO_VERSION=${$1}
        shift;;
    esac
done

if [[ ! -d $OUTPUT_FOLDER ]] ; then
    mkdir -p $OUTPUT_FOLDER
fi

docker run --name build-upgrade-v${GO_VERSION} --privileged --rm -v ${WORKSPACE}:/workspace -v ${OUTPUT_FOLDER}:/output --workdir /workspace golang:${GO_VERSION} "bash build/buildInsideContainer.sh"