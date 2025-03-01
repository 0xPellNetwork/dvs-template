#!/usr/bin/env bash

# Install the required protoc execution tools.
# go install github.com/regen-network/cosmos-proto/protoc-gen-gocosmos
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@latest

# Install goimports for formatting proto go imports.
# go install golang.org/x/tools/cmd/goimports@latest

PROJECT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd)
cd $PROJECT_DIR


go install github.com/cosmos/gogoproto/protoc-gen-gocosmos
go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.16.0

# Generate Gogo proto code.
cd proto
rm -rf generate
buf generate

# Move proto files to the right places.
cd ..
cp -r proto/generate/cosmos/example/dvs/* ./dvs/

rm -rf proto/generate

# # Format proto go imports.
# # goimports -w .