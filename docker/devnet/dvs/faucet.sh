#!/usr/bin/env bash
set -e
set -x

# faucet_to_key <receiver_key> <amount>,
# will send amount of eth to <RECEIVER_KEY> from <ROOT_KEY> in the <ETH_RPC_URL> network
function faucet_to_key {
  RECEIVER_ADDRESS=$(cast wallet address --private-key "${RECEIVER_KEY}")
  ROOT_ADDRESS=$(cast wallet address --private-key "${ROOT_KEY}")
  AMOUNT=$(printf "%0.f" "${2:-1e18}")

  ## By default, cast will use $ETH_RPC_URL environment variable as the RPC URL
  ROOT_BALANCE=$(cast balance "$ROOT_ADDRESS" --rpc-url "$ETH_RPC_URL")
  echo "Root balance: $ROOT_BALANCE"

  ## If cast send throws an error like "duplicate field",
  ## please update the version of forge of eth container to the latest version
  cast send "$RECEIVER_ADDRESS" --value "$AMOUNT" --private-key "$ROOT_KEY" --rpc-url "$ETH_RPC_URL"
  RECEIVER_BALANCE=$(cast balance "$RECEIVER_ADDRESS" --rpc-url "$ETH_RPC_URL")
  echo "Receiver balance: $RECEIVER_BALANCE"
}

if [ -z "$ROOT_KEY" ]; then
  echo "ROOT_KEY is not set"
  exit 1
fi

if [ -z "$RECEIVER_KEY" ]; then
  echo "RECEIVER_KEY is not set"
  exit 1
fi

faucet_to_key
