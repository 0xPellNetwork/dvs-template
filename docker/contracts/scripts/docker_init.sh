#!/bin/bash
set -x
set -e

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function load_defaults {
  export APP_DIR="/app/dvs-contracts-template"
  export CONTRACTS_PATH="/app/dvs-contracts-template/lib/pell-middleware-contracts/lib/pell-contracts"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
}

function eth_healthcheck {
  while true; do
    cast block-number --rpc-url $ETH_RPC_URL
    if [ $? -eq 0 ]; then
      echo "Eth node is ready"
      break
    fi
    echo "Eth node is not ready, retrying in 1 second..."
    sleep 1
  done
}

function deploy_contracts {
  start_time=$(date +%s)
  cd $APP_DIR

  # if deployments folder is empty, deploy contracts
  if [ ! -d "./deployments/localhost" ]; then
    echo "Deploying contracts"
    rm -rf ./deployments/*

    # deploy pell evm
    cd $CONTRACTS_PATH
    npx hardhat deploy --deploy-scripts deploy_restaking --network localhost
    npx hardhat deploy --deploy-scripts deploy_pell --network localhost
    npx hardhat deploy --deploy-scripts deploy_service_omni --network localhost
    npx hardhat update-delegation-connector --network localhost

    # deploy incredible squaring
    cd $APP_DIR
    npx hardhat --network localhost deploy
  else
    echo "Contracts already deployed"
  fi

  # for healthcheck
  touch /root/contracts_deployed_completed
  echo "Total deployment time: $(($(date +%s) - start_time)) seconds"
}

# listen pell events
function listen_pell_events {
  cd $CONTRACTS_PATH
  npx hardhat --network localhost listen-pell-events
}

# start sshd
/usr/sbin/sshd &

logt "Load Default Values for ENV Vars if not set."
load_defaults

logt "Wait for eth to be ready"
eth_healthcheck

logt "Deploy Contracts"
deploy_contracts

logt "Listen Pell Events"
listen_pell_events
