# dvs-template Setup and Usage Guide

In this template, we demonstrate how to build a square number calculation as an example to guide you through the development process. This example is representative of most Layer 2 applications, featuring off-chain computation and writing the results back to the blockchain.

## 1. Install Go Environment

First, ensure that the Go programming language is installed on your system. You can follow the official installation guide here: [Go Official Installation Guide](https://golang.org/doc/install)

## 2. Clone the Repository and Install Dependencies

Use the following commands to clone the project repository and install the necessary Go dependencies:

```bash
git clone https://github.com/0xPellNetwork/dvs-template.git
cd dvs-template && go mod tidy
```

## 3. Modify and Generate Proto Files

### Proto File Definitions

- **Proto Message Definition**: [`/proto/dvs/squared/task.proto`](proto/dvs/squared/task.proto)
- **Proto Service for DVS Request Definition**: [`/proto/dvs/squared/dvs_request.proto`](proto/dvs/squared/dvs_request.proto)
- **Proto Service for DVS Response Definition**: [`/proto/dvs/squared/dvs_response.proto`](proto/dvs/squared/dvs_response.proto)

After modifying the proto files, generate the Go files using the following command:

```bash
make proto
```

## 4. Modify DVS Module

### 1. Customize Result Handler Interface Implementation

- **File Path**: [`dvs/squared/result/handler.go`](dvs/squared/result/handler.go)
- **Description**: Implement the result handler interface based on the DVS request message.

### 2. Customize DVS Server Handler Interface Implementation

- **Request Handler**: [`/dvs/squared/server/request.go`](dvs/squared/server/request.go)
  - Implement the DVS server handler interface based on `dvs_request.proto`.
  
- **Response Handler**: [`/dvs/squared/server/response.go`](dvs/squared/server/response.go)
  - Implement the DVS server handler interface based on `dvs_response.proto`.

### 3. Modify DVS Module Types

- **File Path**: [`/dvs/squared/types/module.go`](dvs/squared/types/module.go)
- **Description**: Modify the dvs module name

## 5. Run Tests

### Unit Tests

Run unit tests using the following command:

```bash
make test
```

### Integration Tests

Run integration tests using the following commands:

```bash
make docker-build && make docker-test
```

## 6. Build Executable

Build the project's executable using the following command:

```bash
make build
```
