#!/bin/bash

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

  export CHAIN_ID_BSC=${CHAIN_ID_BSC:-97}
  export CHAIN_ID_SEPOLIA=${CHAIN_ID_SEPOLIA:-11155111}
  export SERVICE_CHAIN_RPC_URL_BSC=${SERVICE_CHAIN_RPC_URL_BSC:-https://bsc-testnet.public.blastapi.io}
  export SERVICE_CHAIN_RPC_URL_SEPOLIA=${SERVICE_CHAIN_RPC_URL_SEPOLIA:-https://eth-sepolia.public.blastapi.io}

  export AGGREGATOR_RPC_PORT=${AGGREGATOR_RPC_PORT:-26653}
  export AGGREGATOR_RPC_LADDR=${AGGREGATOR_RPC_LADDR:-0.0.0.0:$AGGREGATOR_RPC_PORT}
  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}

  export AGGREGATOR_INDEXER_START_HEIGHT=${AGGREGATOR_INDEXER_START_HEIGHT:-148211}
  export AGGREGATOR_INDEXER_BATCH_SIZE=${AGGREGATOR_INDEXER_BATCH_SIZE:-1000}
  # nanoseconds
  export AGGREGATOR_INDEXER_LISTEN_INTERVAL=${AGGREGATOR_INDEXER_LISTEN_INTERVAL:-5000000000}
}

function set_registry_router_address() {
  # if $REGISTRY_ROUTER_ADDRESS is not set, fetch it from RegistryRouterAddress.json
  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    # TODO(kevin): should get address from contract
    # TODO(@jimmy): seems it can'be retrieved from contract, need to check
    export REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)
  else
    echo "Using provided REGISTRY_ROUTER_ADDRESS: $REGISTRY_ROUTER_ADDRESS"
  fi
}

function init_aggregator {
  pelldvs init --home $PELLDVS_HOME

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" ~/.pelldvs/config/config.toml
  }
  local PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")

  update-config rpc_url "$ETH_RPC_URL"
  update-config pell_registry_router_address "$REGISTRY_ROUTER_ADDRESS"
  update-config pell_delegation_manager_address "$PELL_DELEGATION_MNAGER"

  mkdir -p $PELLDVS_HOME/config

  cat <<EOF > $PELLDVS_HOME/config/aggregator.json
{
    "aggregator_rpc_server": "$AGGREGATOR_RPC_LADDR",
    "operator_response_timeout": "10s",
    "pell_registry_router_address": "$REGISTRY_ROUTER_ADDRESS",
    "chain_config_path": "$PELLDVS_HOME/config/chain.detail.json",
    "indexer_start_height": $AGGREGATOR_INDEXER_START_HEIGHT,
    "indexer_batch_size": $AGGREGATOR_INDEXER_BATCH_SIZE,
    "indexer_listen_interval": $AGGREGATOR_INDEXER_LISTEN_INTERVAL
}
EOF

  HARDHAT_DVS_PATH="deployments/bsc-testnet"
  DVS_OPERATOR_KEY_MANAGER_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorKeyManager-Proxy.json")
  DVS_CENTRAL_SCHEDULER_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/CentralScheduler-Proxy.json")
  DVS_OPERATOR_INFO_PROVIDER_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorInfoProvider.json")
  DVS_OPERATOR_INDEX_MANAGER_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorIndexManager-Proxy.json")
  if [ -z "$DVS_OPERATOR_KEY_MANAGER_BSC" ] || [ -z "$DVS_CENTRAL_SCHEDULER_BSC" ] || [ -z "$DVS_OPERATOR_INFO_PROVIDER_BSC" ] || [ -z "$DVS_OPERATOR_INDEX_MANAGER_BSC" ]; then
   echo "Failed to fetch DVS addresses from $HARDHAT_DVS_PATH"
   exit 1
  fi

  HARDHAT_DVS_PATH="deployments/sepolia"
  DVS_OPERATOR_KEY_MANAGER_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorKeyManager-Proxy.json")
  DVS_CENTRAL_SCHEDULER_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/CentralScheduler-Proxy.json")
  DVS_OPERATOR_INFO_PROVIDER_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorInfoProvider.json")
  DVS_OPERATOR_INDEX_MANAGER_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/OperatorIndexManager-Proxy.json")
  if [ -z "$DVS_OPERATOR_KEY_MANAGER_SEPOLIA" ] || [ -z "$DVS_CENTRAL_SCHEDULER_SEPOLIA" ] || [ -z "$DVS_OPERATOR_INFO_PROVIDER_SEPOLIA" ] || [ -z "$DVS_OPERATOR_INDEX_MANAGER_SEPOLIA" ]; then
   echo "Failed to fetch DVS addresses from $HARDHAT_DVS_PATH"
   exit 1
  fi


  cat <<EOF > $PELLDVS_HOME/config/chain.detail.json
{
  "$CHAIN_ID_BSC": {
    "rpc_url": "$SERVICE_CHAIN_RPC_URL_BSC",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER_BSC",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER_BSC",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER_BSC",
    "operator_index_manager_address": "$DVS_OPERATOR_INDEX_MANAGER_BSC"
  },
  "$CHAIN_ID_SEPOLIA": {
    "rpc_url": "$SERVICE_CHAIN_RPC_URL_SEPOLIA",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER_SEPOLIA",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER_SEPOLIA",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER_SEPOLIA",
    "operator_index_manager_address": "$DVS_OPERATOR_INDEX_MANAGER_SEPOLIA"
  }
}
EOF
}

function start_aggregator {
  pelldvs start-aggregator --home "$PELLDVS_HOME"
}

logt "Load Default Values for ENV Vars if not set."
load_defaults
set_registry_router_address

if [ ! -f /root/aggregator_initialized ]; then
  logt "Init aggregator"
  init_aggregator
  touch /root/aggregator_initialized
fi

logt "Starting aggregator..."
start_aggregator