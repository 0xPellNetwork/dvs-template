
set -x
set -e

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export AGGREGATOR_RPC_URL=${AGGREGATOR_RPC_URL:-dvs:26653}
  export OPERATOR_NODE_NAME=${OPERATOR_NODE_NAME:-operator01}

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}
  export OPERATOR_ADDRESS=0x$(pelldvs keys show $OPERATOR_KEY_NAME --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}

## TODO: move operator config to seperated location
function init_pelldvs_config {
  pelldvs init --home $PELLDVS_HOME
  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }

  # update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(fetch_pell_address "registry_router_factory")
  PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")
  PELL_DVS_DIRECTORY=$(fetch_pell_address "dvs_directory_proxy")

  # check values
  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL is not set"
    exit 1
  fi
  if [ -z "$REGISTRY_ROUTER_FACTORY_ADDRESS" ]; then
    echo "REGISTRY_ROUTER_FACTORY_ADDRESS is not set"
    exit 1
  fi
  if [ -z "$PELL_DELEGATION_MNAGER" ]; then
    echo "PELL_DELEGATION_MNAGER is not set"
    exit 1
  fi
  if [ -z "$PELL_DVS_DIRECTORY" ]; then
    echo "PELL_DVS_DIRECTORY is not set"
    exit 1
  fi
  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    echo "REGISTRY_ROUTER_ADDRESS is not set"
    exit 1
  fi

  update-config rpc_url "$ETH_RPC_URL"
  update-config pell_registry_router_factory_address "$REGISTRY_ROUTER_FACTORY_ADDRESS"
  update-config pell_delegation_manager_address "$PELL_DELEGATION_MNAGER"
  update-config pell_dvs_directory_address "$PELL_DVS_DIRECTORY"
  update-config pell_registry_router_address "$REGISTRY_ROUTER_ADDRESS"
  update-config aggregator_rpc_url "$AGGREGATOR_RPC_URL"

  ## FIXME: don't use absolute path for key
  update-config operator_bls_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.bls.key.json"
  update-config operator_ecdsa_private_key_store_path "$PELLDVS_HOME/keys/$OPERATOR_KEY_NAME.ecdsa.key.json"
}

function register_operator {
  OPERATOR_METADATA_URI=https://raw.githubusercontent.com/matthew7251/Metadata/main/Matthew_Metadata.json
  pelldvs client operator register-operator \
    --home $PELLDVS_HOME \
    --from $OPERATOR_KEY_NAME \
    --metadata-uri $OPERATOR_METADATA_URI

  show_operator_registered "$OPERATOR_ADDRESS"
}

load_defaults
init_pelldvs_config

if [ "$1" == "register" ]; then 
  register_operator
else
  show_operator_registered "$OPERATOR_ADDRESS"
fi
