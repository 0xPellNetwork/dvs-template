
set -x

function load_defaults {
  export ADMIN_KEY=${ADMIN_KEY}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
  export OPERATOR_KEY_NAME=${OPERATOR_KEY_NAME:-operator}
}

function remove_and_import_operator_key {
  rm -rf "$PELLDVS_HOME"/keys/${OPERATOR_KEY_NAME}.ecdsa.key.json
  echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure ${OPERATOR_KEY_NAME} $OPERATOR_KEY --home $PELLDVS_HOME >/dev/null

  export OPERATOR_ADDRESS=$(pelldvs keys show ${OPERATOR_KEY_NAME} --home $PELLDVS_HOME | awk '/Key content:/{getline; print}' | head -n 1 | jq -r .address)

  ## To register operator in the DVS, we need the operator's BLS key with the same name
  rm -rf "$PELLDVS_HOME"/keys/${OPERATOR_KEY_NAME}.bls.key.json
  echo -ne '\n\n' | pelldvs keys import --key-type bls --insecure ${OPERATOR_KEY_NAME} $OPERATOR_BLS_KEY --home $PELLDVS_HOME >/dev/null

}

function check_envs() {
    if [ -z "$OPERATOR_KEY" ]; then
        echo "OPERATOR_KEY is not set"
        exit 1
    fi

    if [ -z "$OPERATOR_BLS_KEY" ]; then
        echo "OPERATOR_BLS_KEY is not set"
        exit 1
    fi

    if [ -z "$PELLDVS_HOME" ]; then
        echo "PELLDVS_HOME is not set"
        exit 1
    fi

    if [ -z "$OPERATOR_KEY_NAME" ]; then
        echo "OPERATOR_KEY_NAME is not set"
        exit 1
    fi
}


load_defaults
check_envs

remove_and_import_operator_key

echo "import_operator_key_forcibly done"
