
set -x

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/dvs-contracts-template/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/dvs-contracts-template/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
}

function task_gateway_healthcheck {
  set +e
  local container_name=$1
  while true; do
    ssh $container_name "test -f /root/operator_initialized"
    if [ $? -eq 0 ]; then
      echo "Operator initialized, proceeding to the next step..."
      break
    fi
    echo "Operator not initialized, retrying in 2 second..."
    sleep 2
  done
  ## Wait for operator to be ready
  sleep 3
  set -e
}

function assert_eq {
  if [ "$1" != "$2" ]; then
    echo "[FAIL] Expected $1 to be equal to $2"
    exit 1
  fi
  echo "[PASS] Expected $1 to be equal to $2"
}

load_defaults
task_gateway_healthcheck operator
