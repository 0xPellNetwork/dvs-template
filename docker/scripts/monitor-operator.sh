#!/usr/bin/env bash

OPERATOR_NODE_LIST=(
  "01"
  "02"
)
LOG_FILE="${LOG_FILE:-$HOME/.pelldvs-homes/logs/monitor-operator.log}"
mkdir -p ~/.pelldvs-homes/logs

logtf() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1" | tee -a $LOG_FILE
}

need_restart=0
reason=""

function check_if_websocket_lost() {
    local node_id=$1
    if [ -z "$node_id" ]; then
        logtf "node_id is not set"
        exit 1
    fi
    local container_name="pelldvs-example-testnet-operator$node_id-1"
    # if docker logs -n 1 $container_name | grep -q "websocket: close 1006 (abnormal closure): unexpected EOF"; then
    # set container shoule be restart and set reason to "websocket: close 1006 (abnormal closure): unexpected EOF"

    if docker logs -n 1 $container_name | grep -q "websocket: close 1006 (abnormal closure): unexpected EOF"; then
        # save log to file
        need_restart=1
        reason="websocket: close 1006 (abnormal closure): unexpected EOF"
    fi

    if docker logs -n 1 $container_name | grep -q "websocket: close 1013: Connection timeout exceeded"; then
        # save log to file
        need_restart=1
        reason="websocket: close 1013: Connection timeout exceeded"
    fi

    # check for: websocket: close 1001
    if docker logs -n 1 $container_name | grep -q "websocket: close 1001"; then
        # save log to file
        need_restart=1
        reason="websocket: close 1001"
    fi

}

function check_if_container_running() {
 local node_id=$1
    if [ -z "$node_id" ]; then
        logtf "node_id is not set"
        exit 1
    fi
    local container_name="pelldvs-example-testnet-operator${node_id}-1"
    if ! docker ps | grep -q $container_name; then
        need_restart=1
        reason="container is not running"
    fi
}

function check_service() {
    local node_id=$1
    if [ -z "$node_id" ]; then
        logtf "node_id is not set"
        exit 1
    fi

    check_if_websocket_lost $node_id
    check_if_container_running $node_id

 if [ $need_restart -eq 1 ]; then
        logtf "is_need Need restart of operator${node_id} for: $reason"

        if [ "$reason" == "container is not running" ]; then
          logtf "Last 5 lines of the log for operator${node_id}:"
          docker logs -n 5 "pelldvs-example-testnet-operator${node_id}-1" | tee -a $LOG_FILE
        fi

        docker restart "pelldvs-example-testnet-operator${node_id}-1"
        logtf "has_restarted operator${node_id}"
    else
        logtf "no_need No need to restart operator${node_id}"
    fi

    need_restart=0
    reason=""

}

for node_id in "${OPERATOR_NODE_LIST[@]}"; do
  check_service $node_id
done
