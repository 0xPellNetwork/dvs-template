#!/bin/bash

set -e

VERSION=${VERSION:-$(grep 'github.com/squaredLabs/iwasm' go.mod | grep -o 'v[0-9.]\+')}
PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
## Download the lib
DOWNLOAD_URL="https://github.com/squaredLabs/iwasm/raw/refs/tags/${VERSION}/testutils/data/processor.wasm"
INSTALL_DIR="${PROJECT_DIR}/tests/iwasm/testdata"

if [ ! -f "${INSTALL_DIR}/processor.wasm" ]; then
    echo "Downloading wasm testdata from ${DOWNLOAD_URL}"
    HTTP_STATUS=$(curl -L \
        -w "%{http_code}" \
        -o "${INSTALL_DIR}/processor.wasm" \
        "${DOWNLOAD_URL}")

    if [ "$HTTP_STATUS" -ne 200 ]; then
        echo "Error: Failed to download wasm testdata. HTTP status: ${HTTP_STATUS}"
        exit 1
    fi
else
    echo "Wasm testdata already exists at ${INSTALL_DIR}/processor.wasm"
fi