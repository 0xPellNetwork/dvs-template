

source "$(dirname "$0")/utils.sh"

function load_defaults {
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export OPERATOR_ADDRESS=$(pelldvs keys show $OPERATOR_KEY_NAME --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)
}

load_defaults
faucet $OPERATOR_ADDRESS