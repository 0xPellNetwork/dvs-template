#!/usr/bin/env bash

set -e
set -x

function load_defaults {
  export GITHUB_TOKEN=${GITHUB_TOKEN}
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}
}

fetch_contract_address() {
  curl -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3.raw" \
      https://api.github.com/repos/0xPellNetwork/contracts/contents/$1 | jq -r '.address'
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

function update_pelldvs_config {
  pelldvs init --home "$PELLDVS_HOME"
  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }
  update-config rpc_url "$ETH_RPC_URL"
  update-config registry_router_address "$REGISTRY_ROUTER_ADDRESS"
}

function setup_group_config_bsctestnet() {
  STBTC_STRATEGY_ADDRESS=$(fetch_contract_address "deployments/$NETWORK/stBTC-Strategy-Proxy.json")
  BTCB_STRATEGY_ADDRESS=$(fetch_contract_address "deployments/$NETWORK/BTCB-Strategy-Proxy.json")

  # if STBTC strategy address is not set, exit
  if [ -z "$STBTC_STRATEGY_ADDRESS" ]; then
    echo "stBTC strategy address is not set"
    exit 1
  fi
  # if BTCB strategy address is not set, exit
  if [ -z "$BTCB_STRATEGY_ADDRESS" ]; then
    echo "BTCB strategy address is not set"
    exit 1
  fi

  cat > $PELLDVS_HOME/group-0-config.json <<EOF
{
  "minimum_stake": 0,
  "pool_params": [
    {
      "chain_id": $CHAIN_ID,
      "multiplier": 1,
      "pool": "$STBTC_STRATEGY_ADDRESS"
    },
    {
      "chain_id": $CHAIN_ID,
      "multiplier": 1,
      "pool": "$BTCB_STRATEGY_ADDRESS"
    }
  ],
  "operator_set_params": {
    "kick_bi_ps_of_operator_stake": 10,
    "kick_bi_ps_of_total_stake": 10,
    "max_operator_count": 1000
  }
}
EOF
}

function setup_group_config_sepolia() {
  WBTC_STRATEGY_ADDRESS=$(fetch_contract_address "deployments/sepolia/WBTC-Strategy-Proxy.json")
  # if wbtc strategy address is not set, exit
  if [ -z "$WBTC_STRATEGY_ADDRESS" ]; then
    echo "WBTC strategy address is not set"
    exit 1
  fi

  cat > $PELLDVS_HOME/group-0-config.json <<EOF
{
  "minimum_stake": 0,
  "pool_params": [
    {
      "chain_id": $CHAIN_ID,
      "multiplier": 1,
      "pool": "$WBTC_STRATEGY_ADDRESS"
    }
  ],
  "operator_set_params": {
    "kick_bi_ps_of_operator_stake": 10,
    "kick_bi_ps_of_total_stake": 10,
    "max_operator_count": 1000
  }
}
EOF
}

function setup_group_config() {
    network=$1
    case $network in
    "bsc-testnet")
      setup_group_config_bsctestnet
      ;;
    "sepolia")
      setup_group_config_sepolia
      ;;
    *)
      echo "Unsupported network: $network"
      exit 1
      ;;
    esac
}

function create_group {
  setup_group_config $NETWORK

  pelldvs client dvs create-group \
    --home $PELLDVS_HOME \
    --from admin \
    --config $PELLDVS_HOME/group-0-config.json
}

function show_group {
  GROUP_COUNT=$(cast call "$REGISTRY_ROUTER_ADDRESS" "groupCount()" --rpc-url "$ETH_RPC_URL")
  echo "Group Count From Registry Router in Pell EVM: $GROUP_COUNT"
}

function check_envs() {
  if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN is not set. Exiting."
    exit 1
  fi

  if [ -z "$NETWORK" ]; then
    echo "NETWORK is not set. Exiting."
    exit 1
  fi

  if [ -z "$CHAIN_ID" ]; then
    echo "CHAIN_ID is not set. Exiting."
    exit 1
  fi

  if [ -z "$PELLDVS_HOME" ]; then
    echo "PELLDVS_HOME is not set. Exiting."
    exit 1
  fi

  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL is not set. Exiting."
    exit 1
  fi

  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    echo "REGISTRY_ROUTER_ADDRESS is not set. Exiting."
    exit 1
  fi
}

load_defaults
set_registry_router_address

check_envs

if [ "$1" == "create" ]; then
  update_pelldvs_config
  create_group
else
  show_group
fi
