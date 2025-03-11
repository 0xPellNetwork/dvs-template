#!/bin/bash

set -x
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
  export AGGREGATOR_RPC_SERVER=${AGGREGATOR_RPC_SERVER:-dvs:26653}
  export GATEWAY_PORT=${GATEWAY_PORT:-8949}
}

function dvs_healthcheck {
  set +e
  while true; do
    curl -s $AGGREGATOR_RPC_SERVER >/dev/null
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


## FIXME: remove this logic after fix. Operator should never use admin key.
function setup_admin_key {
  export ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
  if ! pelldvs keys show admin --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure admin $ADMIN_KEY --home $PELLDVS_HOME >/dev/null
  fi

  export ADMIN_ADDRESS=$(pelldvs keys show admin --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}



function setup_task_gateway_config {
  setup_admin_key
  SERVICE_MANAGER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json" | jq -r .address)
  mkdir -p $PELLDVS_HOME/config

  cat <<EOF > $PELLDVS_HOME/config/task_gateway.config.json
{
  "server_addr": "0.0.0.0:$GATEWAY_PORT",
  "gateway_key_path": "$PELLDVS_HOME/keys/admin.ecdsa.key.json",
  "chains": {
    "1337": {
      "rpc_url": "$ETH_RPC_URL",
      "contract_address": "$SERVICE_MANAGER_ADDRESS",
      "gas_limit": 1000000,
      "chain_id": 1337
    }
  }
}
EOF
}

function start_task_gateway {
  squaringd start-chain-connector --home "$PELLDVS_HOME"
}

## start sshd
/usr/sbin/sshd

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Check if DVS is ready"
dvs_healthcheck

if [ ! -f /root/task_gatewa_initialized ]; then
  logt "Setup task gateway config"
  setup_task_gateway_config
  touch /root/task_gatewa_initialized
fi

logt "Starting task gateway..."
start_task_gateway
