set -x

function load_defaults {
  export NETWORK=${NETWORK:-bsc-testnet}
  export CHAIN_ID=${CHAIN_ID:-97}
  export HARDHAT_DVS_PATH="deployments/$NETWORK"
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
}

function setup_root_key {
  if [ -z "$ROOT_KEY" ]; then
    echo "ROOT_KEY is not set"
    exit 1
  fi
  if [ -z "$PELLDVS_HOME" ]; then
    echo "PELLDVS_HOME is not set"
    exit 1
  fi

  ## If root key is not imported, import it
  if ! pelldvs keys show root --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    # Root key is the key from Pell Network's testnet used to fund
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure root $ROOT_KEY --home $PELLDVS_HOME >/dev/null
  fi
  export ROOT_ADDRESS=$(pelldvs keys show root --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | jq -r .address)
}

function check_github_token() {
    if [ -z "$GITHUB_TOKEN" ]; then
        echo "GITHUB_TOKEN is not set"
        exit 1
    fi
}

function fetch_dvs_address() {
  check_github_token
  curl -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3.raw" \
      https://api.github.com/repos/0xPellNetwork/dvs-contracts-template/contents/$1 | jq -r '.address'
}

function fetch_staking_address() {
  check_github_token
  curl -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3.raw" \
      https://api.github.com/repos/0xPellNetwork/contracts/contents/$1 | jq -r '.address'
}

function fetch_pell_address {
  KEY=$1
  curl https://raw.githubusercontent.com/0xPellNetwork/network-config/refs/heads/main/testnet/system_contract.json | jq -r ".$KEY"
}

function faucet {
  setup_root_key
  RECEIVER_ADDRESS="$1"
  AMOUNT=$(printf "%0.f" "${2:-1e18}")

  if [ -z "$RECEIVER_ADDRESS" ]; then
    echo "RECEIVER_ADDRESS is not set"
    exit 1
  fi

  ## By default, cast will use $ETH_RPC_URL environment variable as the RPC URL
  ROOT_BALANCE=$(cast balance "$ROOT_ADDRESS")
  echo "Root balance: $ROOT_BALANCE"

  ## If cast send throws an error like "duplicate field", 
  ## please update the version of forge of eth container to the latest version
  cast send "$RECEIVER_ADDRESS" --value "$AMOUNT" --private-key "$ROOT_KEY"
  RECEIVER_BALANCE=$(cast balance "$RECEIVER_ADDRESS")
  echo "Receiver balance: $RECEIVER_BALANCE"
}

function show_operator_registered {
  local ADDRESS=$1
  local PELL_DELEGATION_MNAGER=$(fetch_pell_address "delegation_manager_proxy")
  local IS_PELL_OPERATOR=$(cast call $PELL_DELEGATION_MNAGER "isOperator(address)" $ADDRESS)
  echo "Is pell operator: $IS_PELL_OPERATOR"
}

function show_dvs_operator_info {
  local ADDRESS=$1
  if [ -z "$ADDRESS" ]; then
    echo "ADDRESS is not set"
    exit 1
  fi

  DVS_CENTRAL_SCHEDULER=$(fetch_dvs_address "$HARDHAT_DVS_PATH/CentralScheduler-Proxy.json")
  if [ -z "$DVS_CENTRAL_SCHEDULER" ]; then
    echo "DVS_CENTRAL_SCHEDULER is not set"
    exit 1
  fi

  ## Get operator info -> (operator_id, status),
  ## status: 0 -> NEVER, 1 -> REGISTERED, 2 -> DEREGISTERED
  cast call "$DVS_CENTRAL_SCHEDULER" "getOperator(address)((bytes32,uint8))" $ADDRESS --rpc-url $SERVICE_CHAIN_RPC_URL
}


load_defaults
