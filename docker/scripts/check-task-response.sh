#!/usr/bin/env bash

HERE=$(dirname "$0")
PARRENT_DIR=$(dirname "$HERE")

DOCKER_COMPOSE_FILE=${DOCKER_COMPOSE_FILE:-$PARRENT_DIR/docker-compose.yml}
echo "Checking task response... $DOCKER_COMPOSE_FILE"

CHECK_CONTENT="Task response sent successfully"

for i in {1..2}; do
  # check if the task response is sent successfully
  if docker compose -f $DOCKER_COMPOSE_FILE logs operator0$i |grep 'Task response sent successfully'; then
    echo -e "\t\t\t\tTask response sent successfully in operator0$i"
    continue
  else
    echo -e "\t\t\t\tTask response not sent successfully in operator0$i"
    exit 1
  fi
done

echo -e "\n"
echo "Task response sent successfully in all operators"
echo -e "\n"

exit 0
