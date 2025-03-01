#!/bin/bash

set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}

  export CHAIN_ID_BSC=${CHAIN_ID_BSC:-97}
  export CHAIN_ID_SEPOLIA=${CHAIN_ID_SEPOLIA:-11155111}
  export SERVICE_CHAIN_RPC_URL_BSC=${SERVICE_CHAIN_RPC_URL_BSC:-https://bsc-testnet.public.blastapi.io}
  export SERVICE_CHAIN_RPC_URL_SEPOLIA=${SERVICE_CHAIN_RPC_URL_SEPOLIA:-https://eth-sepolia.public.blastapi.io}

  export GATEWAY_PORT=${GATEWAY_PORT:-8949}
  export GATEWAY_KEY=${GATEWAY_KEY}
  export AGGREGATOR_RPC_URL=${AGGREGATOR_RPC_URL:-dvs:26653}
}

function dvs_healthcheck {
  set +e
  while true; do
    curl -s $AGGREGATOR_RPC_URL >/dev/null
    if [ $? -eq 52 ]; then
      echo "DVS RPC port is ready, proceeding to the next step..."
      break
    fi
    echo "DVS RPC port not ready, retrying in 2 seconds..."
    sleep 2
  done
  ## Wait for aggregator to be ready
  sleep 3
  set -e
}

function setup_gateway_key {
  if [ -z "$GATEWAY_KEY" ]; then
    echo "GATEWAY_KEY is not set. Exiting."
    exit 1
  fi
  
  if ! pelldvs keys show gateway --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure gateway $GATEWAY_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export GATEWAY_ADDRESS=$(pelldvs keys show gateway --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}

function setup_gateway_config {
  setup_gateway_key
  HARDHAT_DVS_PATH="deployments/bsc-testnet"
  INCREDIBLE_SQUARING_SERVICE_MANAGER_ADDRESS_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json")
  HARDHAT_DVS_PATH="deployments/sepolia"
  INCREDIBLE_SQUARING_SERVICE_MANAGER_ADDRESS_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json")

  cat <<EOF > $PELLDVS_HOME/config/task_gateway.config.json
{
  "server_addr": "0.0.0.0:$GATEWAY_PORT",
  "gateway_key_path": "$PELLDVS_HOME/keys/gateway.ecdsa.key.json",
  "chains": {
    "$CHAIN_ID_BSC": {
      "rpc_url": "$SERVICE_CHAIN_RPC_URL_BSC",
      "contract_address": "$INCREDIBLE_SQUARING_SERVICE_MANAGER_ADDRESS_BSC",
      "gas_limit": 1000000,
      "chain_id": $CHAIN_ID_BSC
    },
    "$CHAIN_ID_SEPOLIA": {
      "rpc_url": "$SERVICE_CHAIN_RPC_URL_SEPOLIA",
      "contract_address": "$INCREDIBLE_SQUARING_SERVICE_MANAGER_ADDRESS_SEPOLIA",
      "gas_limit": 1000000,
      "chain_id": $CHAIN_ID_SEPOLIA
    }
  }
}
EOF
}

function start_gateway {
  squaringd start-task-gateway --home $PELLDVS_HOME
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

#logt "Check if DVS is ready"
#dvs_healthcheck

# if GATEWAY_KEY is not set, exit
if [ -z "$GATEWAY_KEY" ]; then
  echo "GATEWAY_KEY is not set. Exiting."
  exit 1
fi

logt "Setup gateway config"
setup_gateway_config

touch /root/gateway_initialized

logt "Starting gateway..."
start_gateway
