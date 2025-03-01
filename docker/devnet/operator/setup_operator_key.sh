
set -x

function load_defaults {
  export ADMIN_KEY=${ADMIN_KEY}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
}

function setup_operator_key {
  # check if file "$PELLDVS_HOME"/keys/${OPERATOR_KEY_NAME}.ecdsa.key.json exists
  # if not, import or create a new key
  if [ ! -f "$PELLDVS_HOME"/keys/${OPERATOR_KEY_NAME}.ecdsa.key.json ]; then
    if [ -z "$OPERATOR_KEY" ]; then
      echo  -ne '\n\n' | pelldvs keys create ${OPERATOR_KEY_NAME} --key-type=ecdsa --insecure > /tmp/operator.key
    else
      echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure ${OPERATOR_KEY_NAME} $OPERATOR_KEY --home $PELLDVS_HOME >/dev/null
    fi
  fi

  export OPERATOR_ADDRESS=$(pelldvs keys show ${OPERATOR_KEY_NAME} --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)

  ## To register operator in the DVS, we need the operator's BLS key with the same name
  if [ ! -f "$PELLDVS_HOME"/keys/${OPERATOR_KEY_NAME}.bls.key.json ]; then
    # if $OPERATOR_BLS_KEY is not set, create a new key, otherwise import the key
    if [ -z "$OPERATOR_BLS_KEY" ]; then
      echo  -ne '\n\n' | pelldvs keys create ${OPERATOR_KEY_NAME} --key-type=bls --insecure > /tmp/operator_bls.key
    else
      echo -ne '\n\n' | pelldvs keys import --key-type bls --insecure ${OPERATOR_KEY_NAME} $OPERATOR_BLS_KEY --home $PELLDVS_HOME >/dev/null
    fi
  fi

}

load_defaults
setup_operator_key

echo "setup_operator_key done"
