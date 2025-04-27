
set -x

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/dvs-contracts-template/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/dvs-contracts-template/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
  export QUERY_SERVER_ADDR=${QUERY_SERVER_ADDR:-operator:8123}
  export OPERATOR_SEVER_NAME_LIST=${OPERATOR_SEVER_NAME_LIST:-"operator"}
}

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function operator_healthcheck {
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

iterate_operators() {
    for operator in "$OPERATOR_SEVER_NAME_LIST"; do
        operator_healthcheck "$operator"
    done
}

load_defaults
iterate_operators $OPERATOR_SEVER_NAME_LIST

asyncURL="http://operator:26657/request_dvs_async?data=%22221111111111=57945678901234567890123456789017%22&height=111&chainid=1337"
asyncResponse=$(curl -sS -H "Accept: application/json" -X GET "$asyncURL")

asyncExpectedStr="{\"jsonrpc\":\"2.0\",\"id\":-1,\"result\":{\"hash\":\"629E74E9238DC7C66902734AD02BA77B3EA52AF68B352190012780537F220516\"}}"

if echo "$asyncResponse" | grep -q "$asyncExpectedStr"; then
	echo "$asyncResponse"
	echo "test async rpc task: ok"
else
	exit 1
fi


logt "test done"
