# dvs-template Setup and Usage Guide

## 1. Build the Docker image

build the Docker image using the following command:

```bash
make docker-build
# cd docker && docker compose -f docker-compose.build.yml build
```

This primarily builds the following three images:

- hardhat: Contract environment for the example, responsible for deploying on-chain contracts to the EVM environment.
- pelldvs: Binary environment for pelldvs, responsible for starting the aggregator and managing the DVS.
- operator: Off-chain runtime environment for the example. The operator performs off-chain computations and returns squared values for numbers requested by on-chain contracts.

## 2. Start Docker Deployment

The example environment primarily includes the following five services (start them in order; the full deployment takes about 2 minutes):

- eth: EVM environment launched using Anvil.
- hardhat: Deploys the Pell protocol's contracts on-chain, including staking EVM, Pell EVM, and service EVM contracts.
- emulator: A simulator for the Pell blockchain that synchronizes staking states across the staking EVM, Pell EVM, and service EVM.
- dvs: Managed by the DVS developers, providing signature aggregation services (Aggregator) and submitting final results on-chain (via TaskGateway).
- operator: Configures and starts node services. The workflow includes:
  - Listening for contract tasks on the EVM (eth).
  - Processing tasks, such as squaring numbers.
  - Using the DVS signature aggregation service to obtain results signed by all operators.
  - Sending the result to the DVS's TaskGateway for on-chain submission.
- test: Sends tasks to test the entire process.

To start the full service suite, use the following command:

```bash
make docker-up-all
# Equivalent command:
# cd docker && docker compose down -v && docker compose up -d
```

After about 2 minutes, all services should be fully started. You can check the test logs to ensure everything is running smoothly:

```bash
make docker-test
```

This task will send a request (with the number 2) and is expected to return the result 4. The expected output is as follows:

```bash
test-1  | ++ cast call 0x2B0d36FACD61B71CC05ab8F3D2355ec3631C0dd5 'taskNumber()(uint32)' --private-key 0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
test-1  | ++ xargs printf %d
test-1  | + TASK_NUMBER=1
test-1  | ++ cast call 0x2B0d36FACD61B71CC05ab8F3D2355ec3631C0dd5 'numberSquareds(uint32)(uint256)' 0
test-1  | + RESULT=4
test-1  | + assert_eq 4 4
test-1  | + '[' 4 '!=' 4 ']'
test-1  | + echo '[PASS] Expected 4 to be equal to 4'
test-1  | [PASS] Expected 4 to be equal to 4
```

# dvs-template
