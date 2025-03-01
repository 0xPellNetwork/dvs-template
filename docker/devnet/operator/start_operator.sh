
#!/bin/bash

set -x
set -e

source "$(dirname "$0")/utils.sh"

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export GATEWAY_ADDR=${GATEWAY_ADDR:-gateway:8949}
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}

  export AGGREGATOR_RPC_URL=${AGGREGATOR_RPC_URL:-dvs:26653}

  export SERVICE_CHAIN_RPC_URL=${SERVICE_CHAIN_RPC_URL:-https://bsc-testnet.blockpi.network/v1/rpc/public}
  export SERVICE_CHAIN_WS_URL=${SERVICE_CHAIN_WS_URL}

  export CHAIN_ID_BSC=${CHAIN_ID_BSC:-97}
  export CHAIN_ID_SEPOLIA=${CHAIN_ID_SEPOLIA:-11155111}
  export SERVICE_CHAIN_RPC_URL_BSC=${SERVICE_CHAIN_RPC_URL_BSC:-https://bsc-testnet.public.blastapi.io}
  export SERVICE_CHAIN_RPC_URL_SEPOLIA=${SERVICE_CHAIN_RPC_URL_SEPOLIA:-https://eth-sepolia.public.blastapi.io}
  export SERVICE_CHAIN_WS_URL_BSC=${SERVICE_CHAIN_WS_URL_BSC:-wss://bsc-testnet-rpc.publicnode.com}
  export SERVICE_CHAIN_WS_URL_SEPOLIA=${SERVICE_CHAIN_WS_URL_SEPOLIA:-wss://ethereum-sepolia-rpc.publicnode.com}

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

function gateway_healthcheck {
  set +e
  while true; do
    curl -s $GATEWAY_ADDR >/dev/null
    if [ $? -eq 52 ]; then
      echo "Gateway is ready, proceeding to the next step..."
      break
    fi
    echo "Gateway not ready, retrying in 2 seconds..."
    sleep 2
  done
  ## Wait for aggregator to be ready
  sleep 3
  set -e
}

## TODO: move operator config to seperated location
function init_pelldvs_config {
  pelldvs init --home $PELLDVS_HOME
  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }

  ## update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(fetch_pell_address "registry_router_factory")
  PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")
  PELL_DVS_DIRECTORY=$(fetch_pell_address "dvs_directory")

  update-config rpc_url "$ETH_RPC_URL"
  update-config pell_registry_router_factory_address "$REGISTRY_ROUTER_FACTORY_ADDRESS"
  update-config pell_delegation_manager_address "$PELL_DELEGATION_MNAGER"
  update-config pell_dvs_directory_address "$PELL_DVS_DIRECTORY"
  update-config pell_registry_router_address "$REGISTRY_ROUTER_ADDRESS"
  update-config aggregator_rpc_url "$AGGREGATOR_RPC_URL"

  ## FIXME: don't use absolute path for key
  update-config operator_bls_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.bls.key.json"
  update-config operator_ecdsa_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.ecdsa.key.json"


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
    "rpc_url": "$SERVICE_CHAIN_WS_URL_BSC",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER_BSC",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER_BSC",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER_BSC",
    "operator_index_manager_address": "$DVS_OPERATOR_INDEX_MANAGER_BSC"
  },
  "$CHAIN_ID_SEPOLIA": {
    "rpc_url": "$SERVICE_CHAIN_WS_URL_SEPOLIA",
    "operator_info_provider_address": "$DVS_OPERATOR_INFO_PROVIDER_SEPOLIA",
    "operator_key_manager_address": "$DVS_OPERATOR_KEY_MANAGER_SEPOLIA",
    "central_scheduler_address": "$DVS_CENTRAL_SCHEDULER_SEPOLIA",
    "operator_index_manager_address": "$DVS_OPERATOR_INDEX_MANAGER_SEPOLIA"
  }
}
EOF
}

function setup_operator_config {
  HARDHAT_DVS_PATH="deployments/bsc-testnet"
  SERVICE_MANAGER_ADDRESS_BSC=$(fetch_dvs_address "$HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json")

  HARDHAT_DVS_PATH="deployments/sepolia"
  SERVICE_MANAGER_ADDRESS_SEPOLIA=$(fetch_dvs_address "$HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json")

  cat <<EOF > $PELLDVS_HOME/config/squaring.config.json
{
  "chain_service_manager_address": {
    "$CHAIN_ID_BSC": "$SERVICE_MANAGER_ADDRESS_BSC",
    "$CHAIN_ID_SEPOLIA": "$SERVICE_MANAGER_ADDRESS_SEPOLIA"
  },
  "gateway_rpc_client_url": "$GATEWAY_ADDR"
}
EOF
}

function start_operator {
  squaringd start-operator --home $PELLDVS_HOME
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "setup operator key"
source "$(dirname "$0")/setup_operator_key.sh"

logt "Check if DVS is ready"
dvs_healthcheck

logt "Check if Gateway is ready"
gateway_healthcheck

logt "Setup operator config"
init_pelldvs_config
setup_operator_config

logt "Starting operator..."
start_operator
