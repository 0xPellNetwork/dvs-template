set -e
set -x

function load_defaults {
  export GITHUB_TOKEN=${GITHUB_TOKEN}
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export SERVICE_CHAIN_RPC_URL=${SERVICE_CHAIN_RPC_URL}
  export REGISTRY_ROUTER_ADDRESS=${REGISTRY_ROUTER_ADDRESS}
}

fetch_contract_address() {
  curl -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3.raw" \
      https://api.github.com/repos/0xPellNetwork/dvs-contracts-template/contents/$1 | jq -r '.address'
}

function update_pelldvs_config {
  pelldvs init --home "$PELLDVS_HOME"

  # if $REGISTRY_ROUTER_ADDRESS is not set, fetch it from RegistryRouterAddress.json
  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    # TODO: should get address from contract
    export REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)
  fi

  update-config() {
    KEY="$1"
    VALUE="$2"
    sed -i "s|${KEY} = \".*\"|${KEY} = \"${VALUE}\"|" $PELLDVS_HOME/config/config.toml
  }
  update-config rpc_url "$ETH_RPC_URL"
  update-config registry_router_address "$REGISTRY_ROUTER_ADDRESS"

}

function add_supported_chain {
  export CENTRAL_SCHEDULER_ADDRESS=$(fetch_contract_address "deployments/$NETWORK/CentralScheduler-Proxy.json")

  if [ -z "$CENTRAL_SCHEDULER_ADDRESS" ]; then
    echo "Central scheduler address is not set"
    exit 1
  fi

  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    echo "REGISTRY_ROUTER_ADDRESS is not set"
    exit 1
  fi
  if [ -z "$SERVICE_CHAIN_RPC_URL" ]; then
    echo "SERVICE_CHAIN_RPC_URL is not set"
    exit 1
  fi
  if [ -z "$CHAIN_ID" ]; then
    echo "CHAIN_ID is not set"
    exit 1
  fi
  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL is not set"
    exit 1
  fi

  pelldvs client dvs register-chain-to-pell \
    --from admin \
    --rpc-url $ETH_RPC_URL \
    --chain-id $CHAIN_ID \
    --registry-router $REGISTRY_ROUTER_ADDRESS \
    --central-scheduler $CENTRAL_SCHEDULER_ADDRESS \
    --dvs-rpc-url $SERVICE_CHAIN_RPC_URL \
    --approver-key-name admin \
    --dvs-from admin
}

function show_supported_chain {
  REGISTRY_ROUTER_ADDRESS=$(cat $PELLDVS_HOME/RegistryRouterAddress.json | jq -r .address)
  cast call $REGISTRY_ROUTER_ADDRESS "groupCount()(uint256)" --rpc-url $ETH_RPC_URL
  cast call $REGISTRY_ROUTER_ADDRESS "getSupportedChain()" --rpc-url $ETH_RPC_URL
  # struct DVSInfo { uint256 chainId; bytes registryCoordinator; bytes ejectionManager; bytes stakeRegistry; }
  cast call $REGISTRY_ROUTER_ADDRESS "supportedChainInfos(uint256)" 0 --rpc-url $ETH_RPC_URL
}

function check_envs() {
  if [ -z "$SERVICE_CHAIN_RPC_URL" ]; then
    echo "SERVICE_CHAIN_RPC_URL is not set"
    exit 1
  fi

  if [ -z "$REGISTRY_ROUTER_ADDRESS" ]; then
    echo "REGISTRY_ROUTER_ADDRESS is not set"
    exit 1
  fi

  if [ -z "$GITHUB_TOKEN" ]; then
    echo "GITHUB_TOKEN is not set"
    exit 1
  fi

  if [ -z "$PELLDVS_HOME" ]; then
    echo "PELLDVS_HOME is not set"
    exit 1
  fi

  if [ -z "$ETH_RPC_URL" ]; then
    echo "ETH_RPC_URL is not set"
    exit 1
  fi

  if [ -z "$NETWORK" ]; then
    echo "NETWORK is not set"
    exit 1
  fi

  if [ -z "$CHAIN_ID" ]; then
    echo "CHAIN_ID is not set"
    exit 1
  fi
}

load_defaults
check_envs

if [ "$1" == "create" ]; then
  update_pelldvs_config
  add_supported_chain
else
  show_supported_chain
fi