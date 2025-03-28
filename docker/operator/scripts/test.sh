
set -x

function load_defaults {
  export HARDHAT_CONTRACTS_PATH="/app/dvs-contracts-template/lib/pell-middleware-contracts/lib/pell-contracts/deployments/localhost"
  export HARDHAT_DVS_PATH="/app/dvs-contracts-template/deployments/localhost"

  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export ETH_RPC_URL=${ETH_RPC_URL:-http://eth:8545}
  export ETH_WS_URL=${ETH_WS_URL:-ws://eth:8545}
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

load_defaults
operator_healthcheck operator01
operator_healthcheck operator02

ADMIN_KEY=0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80
SERVICE_MANAGER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/IncredibleSquaringServiceManager-Proxy.json" | jq -r .address)

## create a new task
NUMBER_TO_BE_SQUARED=$((RANDOM % 10))
RESULT_OF_SQUARED_NUMBER=$((NUMBER_TO_BE_SQUARED * NUMBER_TO_BE_SQUARED))
THREADSHOLD=2 # percentage point
GROUP_NUMBERS=0x00
cast send "$SERVICE_MANAGER_ADDRESS" "createNewTask(uint256,uint32,bytes)" $NUMBER_TO_BE_SQUARED $THREADSHOLD $GROUP_NUMBERS --private-key "$ADMIN_KEY"

## wait for the task to be processed
export TIMEOUT_FOR_TASK_PROCESS=${TIMEOUT_FOR_TASK_PROCESS:-8}
export TIMEOUT_FOR_TASK_PROCESS=$TIMEOUT_FOR_TASK_PROCESS
echo "wait ${TIMEOUT_FOR_TASK_PROCESS} seconds for the task to be processed"
sleep ${TIMEOUT_FOR_TASK_PROCESS}
TASK_NUMBER=$(cast call "$SERVICE_MANAGER_ADDRESS" "taskNumber()(uint32)" --private-key "$ADMIN_KEY" | xargs printf "%d")
RESULT=$(cast call "$SERVICE_MANAGER_ADDRESS" "numberSquareds(uint32)(uint256)" $((TASK_NUMBER - 1)))
assert_eq "$RESULT" "$RESULT_OF_SQUARED_NUMBER"

sleep 5
logt "preparing to check task result from API"
# check via REST API
TASK_ID=$((TASK_NUMBER - 1))
RESPONSE=$(curl -s -X GET "http://operator01:8123/dvs/squared/v1/tasks/${TASK_ID}")
logt "Response from API: $RESPONSE"
RESULT_FROM_API=$(echo $RESPONSE | jq -r .value.result )
logt "Result from API: $RESULT_FROM_API"
assert_eq "$RESULT_FROM_API" "$RESULT_OF_SQUARED_NUMBER"

# cast call "$SERVICE_MANAGER_ADDRESS" "allTaskResponses(uint32)" $((TASK_NUMBER - 1))
# RETRIEVER_ADDRESS=$(ssh hardhat "cat $HARDHAT_DVS_PATH/OperatorStateRetriever.json" | jq -r .address)
# cast call "$RETRIEVER_ADDRESS" "GetGroupsDVSStateAtBlock(uint32)" $TASK_ID --private-key "$ADMIN_KEY"
