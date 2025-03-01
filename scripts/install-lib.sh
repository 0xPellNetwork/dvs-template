#!/bin/bash

set -e
set -x

VERSION=${VERSION:-$(grep 'github.com/squaredLabs/iwasm' go.mod | grep -o 'v[0-9.]\+')}
PLATFORM=${PLATFORM:-$(uname -s | tr '[:upper:]' '[:lower:]')}
ARCH=${ARCH:-$(uname -m)}
PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)

# Map architecture names
case ${ARCH} in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
esac

## Download the lib
LIB_NAME="libruntime-${PLATFORM}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/squaredLabs/iwasm/releases/download/${VERSION}/${LIB_NAME}"
INSTALL_DIR="${PROJECT_DIR}/lib"
mkdir -p "${INSTALL_DIR}"

if [ ! -f "${INSTALL_DIR}/${LIB_NAME}" ]; then
    echo "Downloading libruntime from ${DOWNLOAD_URL}"
    HTTP_STATUS=$(curl -L \
        -w "%{http_code}" \
        -o "${INSTALL_DIR}/${LIB_NAME}" \
        "${DOWNLOAD_URL}")

    if [ "$HTTP_STATUS" -ne 200 ]; then
        echo "Error: Failed to download library. HTTP status: ${HTTP_STATUS}"
        exit 1
    fi
else
    echo "Library ${LIB_NAME} already downloaded, skipping download"
fi

echo "Extracting library ${LIB_NAME} ..."
tar -xzf "${INSTALL_DIR}/${LIB_NAME}" -C "${INSTALL_DIR}" --strip-components=1

echo "Installation completed successfully"
