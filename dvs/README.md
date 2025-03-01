# Minimal DVS Module Example

This document serves as a guide to creating a minimal Decentralized Verification System (DVS) module, which is essential for developers looking to quickly set up a DVS module using the `pelldvs-sdk`.

## Overview

A DVS module is a crucial component in systems that require data verification and integrity checks. This example module is structured to include key directories and files necessary for handling DVS messages and their results.

## Key Components

### 1. Define Protobuf Files

The first step in creating a DVS module is to define the server-related protobuf files. These files are crucial as they specify the protocol messages and gRPC interfaces that the DVS module will use to interact with external systems.

- **Purpose of Protobuf Definitions**:
  - **Protocol Messages**: These are the structured data formats that the DVS module will send and receive. Defining these messages ensures that all parties involved in the communication understand the data structure and semantics.
  - **gRPC Interfaces**: These interfaces define the methods that can be called remotely. By specifying these interfaces, developers can ensure that the DVS module can handle incoming requests and send responses in a standardized manner.

- **Key Components of Protobuf Definitions**:
  - **Task Message Definition**: This defines the structure of the task messages that the DVS module will process.
  - **DVS Request Service Definition**: This service handles incoming requests to the DVS module.
  - **DVS Response Service Definition**: This service manages the responses sent from the DVS module.

- **Steps to Define Protobuf Files**:
  1. **Identify Required Messages**: Determine the data that needs to be exchanged between the DVS module and external systems. This includes request and response messages.
  2. **Define Message Structures**: Use the protobuf syntax to define the structure of each message. This includes specifying fields, data types, and any necessary annotations.
  3. **Specify gRPC Services**: Define the gRPC services that will handle the communication. This involves specifying the service name and the RPC methods, along with their input and output message types.

- **Reference for Definitions**:
  - For specific examples and guidance on defining protobuf files, refer to the existing proto definitions located in the [proto directory](../proto). These examples can provide a template and best practices for structuring your own protobuf files.

By carefully defining the protobuf files, developers can ensure seamless communication and integration of the DVS module with other components in the system. This foundational step is critical for the successful implementation and operation of the DVS module.

### 2. Result Directory

The `result` directory is responsible for implementing the `result handler interface` from the `pelldvs-sdk`. This interface is crucial for processing DVS message results and computing their digests, which are used to verify data integrity.

- **Interface Definition**: The `ResultCustomizedIFace` interface provides two main functions:
  
  ```go
  type ResultCustomizedIFace interface {
      // GetData serializes the result into a byte array
      GetData(proto.Message) ([]byte, error)
      // GetDigest calculates the digest of the result
      GetDigest(proto.Message) ([]byte, error)
  }
  ```

  - `GetData`: This function is responsible for converting the result into a byte array, which is a common format for data transmission and storage.
  - `GetDigest`: This function calculates a digest (a unique hash) of the result, ensuring that the data has not been altered.

### 3. Server Directory

The `server` directory is tasked with implementing both the request and response handlers for DVS messages. These handlers are essential for managing the flow of requests and responses within the DVS module.

- **Request and Response Handling**: Typically defined in protocol buffer (proto) files, these handlers are necessary for the DVS module to function correctly within the `pelldvs` processing flow. They ensure that requests are processed and responses are generated according to the defined protocols.

### 4. module.go

The `module.go` file is the heart of the DVS module, responsible for registering the various components of the module. This includes:

- **Registration Process**:
  - **Request Server Registration**: Ensures that the module can handle incoming requests.
  - **Response Server Registration**: Ensures that the module can send out appropriate responses.
  - **Result Handler Registration**: Integrates the result processing capabilities into the module.

By following these steps, developers can create a fully functional DVS module based on the `pellapp-sdk`. Once the basic structure is in place, developers can focus on implementing the specific logic required for their application.

## Conclusion

This minimal DVS module example provides a foundational structure that can be expanded upon to meet specific data verification needs. By understanding and utilizing the components outlined above, developers can ensure robust data integrity and verification within their systems.
