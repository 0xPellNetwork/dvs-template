#!/bin/bash

set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/dvs-contracts-template/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/dvs-contracts-template/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}

  export AGGREGATOR_RPC_LADDR=${AGGREGATOR_RPC_LADDR:-0.0.0.0:26653}

  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}
}

function set_registry_router_address() {
  # if $REGISTRY_ROUTER_ADDRESS is not set, fetch it from RegistryRouterAddress.json
  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    export REGISTRY_ROUTER_ADDRESS=$(ssh emulator "cat /root/RegistryRouterAddress.json" | jq -r .address)
  else
    echo "Using provided REGISTRY_ROUTER_ADDRESS: $REGISTRY_ROUTER_ADDRESS"
  fi
}

function init_aggregator {
  mkdir -p $PELLDVS_HOME/config
  cat <<EOF > $PELLDVS_HOME/config/aggregator.json
{
    "aggregator_rpc_server": "$AGGREGATOR_RPC_LADDR",
    "operator_response_timeout": "10s",
    "pell_registry_router_address": "$REGISTRY_ROUTER_ADDRESS",
    "chain_config_path": "$PELLDVS_HOME/config/chain.detail.json"
}
EOF

  DVS_OPERATOR_KEY_MANAGER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorKeyManager-Proxy.json" | jq -r .address)
  DVS_CENTRAL_SCHEDULER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/CentralScheduler-Proxy.json" | jq -r .address)
  DVS_OPERATOR_INFO_PROVIDER=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorInfoProvider.json" | jq -r .address)
  cat <<EOF > $PELLDVS_HOME/config/chain.detail.json
{
  "1337": {
    "rpc_url": "$ETH_WS_URL",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER"
  }
}
EOF
}

function start_aggregator {
  echo 'using go run cmd/squaringd/main.go start-aggregator --home "$PELLDVS_HOME"'
  cd /app
  if [ -z "$GITHUB_TOKEN" ]; then echo "GITHUB_TOKEN is not set" && exit 1; fi
  git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/0xPellNetwork".insteadOf "https://github.com/0xPellNetwork"
  go run cmd/squaringd/main.go start-aggregator --home "$PELLDVS_HOME"
}

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Set registry router address"
set_registry_router_address

if [ ! -f /root/aggregator_initialized ]; then
  logt "Init aggregator"
  init_aggregator
  touch /root/aggregator_initialized
fi

logt "Starting aggregator..."
start_aggregator