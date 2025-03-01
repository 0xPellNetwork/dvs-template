set -x

function load_defaults {
  export ADMIN_KEY=${ADMIN_KEY}
  export PELLDVS_HOME=${PELLDVS_HOME:-/root/.pelldvs}
}

function setup_admin_key {
  if ! pelldvs keys show admin --home "$PELLDVS_HOME" >/dev/null 2>&1; then
    echo -ne '\n\n' | pelldvs keys import --key-type ecdsa --insecure admin $ADMIN_KEY --home $PELLDVS_HOME >/dev/null
  fi
}

load_defaults

# if ADMIN_KEY is not set, exit
if [ -z "$ADMIN_KEY" ]; then
  echo "ADMIN_KEY is not set"
  exit 1
fi

setup_admin_key
