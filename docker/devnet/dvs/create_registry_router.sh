set -x
set -e

function load_defaults {
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
}

fetch_pell_address() {
  KEY=$1
  curl https://raw.githubusercontent.com/0xPellNetwork/network-config/refs/heads/main/testnet/system_contract.json | jq -r ".$KEY"
}

function update_pelldvs_config {
  pelldvs init --home "$PELLDVS_HOME"

  ## Update config
  REGISTRY_ROUTER_FACTORY_ADDRESS=$(fetch_pell_address "registry_router_factory")

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }
  update-config rpc_url "$ETH_RPC_URL"
  update-config registry_router_factory_address "$REGISTRY_ROUTER_FACTORY_ADDRESS"
}

function create_registry_router {
  ## Create registry router
  export ADMIN_ADDRESS=$(pelldvs keys show admin --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
  if [ -z "$ADMIN_ADDRESS" ]; then
    echo "Admin address is not set"
    exit 1
  fi

  REGISTRY_ROUTER_ADDRESS_FILE="$PELLDVS_HOME/RegistryRouterAddress.json"
  pelldvs client dvs create-registry-router \
    --home $PELLDVS_HOME \
    --from admin \
    --initial-owner $ADMIN_ADDRESS \
    --dvs-chain-approver $ADMIN_ADDRESS \
    --churn-approver $ADMIN_ADDRESS \
    --ejector $ADMIN_ADDRESS \
    --pauser $ADMIN_ADDRESS \
    --unpauser $ADMIN_ADDRESS \
    --initial-paused-status false \
    --save-to-file $REGISTRY_ROUTER_ADDRESS_FILE \
    --force-save true
}

function check_envs() {
  if [ -z "$PELLDVS_HOME" ]; then
    echo "PELLDVS_HOME is not set"
    exit 1
  fi

  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL is not set"
    exit 1
  fi

  if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN is not set"
    exit 1
  fi
}

load_defaults

check_envs

# if $REGISTRY_ROUTER_ADDRESS_FILE is already created, then we don't need to create it again
if [ ! -f "$REGISTRY_ROUTER_ADDRESS_FILE" ]; then
  update_pelldvs_config
  create_registry_router
else
  echo "Registry router already created"
  # print the address of the registry router
  cat $REGISTRY_ROUTER_ADDRESS_FILE | jq -r .address
fi
