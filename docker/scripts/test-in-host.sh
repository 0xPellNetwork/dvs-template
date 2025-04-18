#!/usr/bin/env bash

HERE=$(dirname "$0")
PARRENT_DIR=$(dirname "$HERE")

cd $PARRENT_DIR;

logt() {
  echo "$(date '+%Y-%m-%d %H:%M:%S') $1"
}

function check_operator_node_ready() {
	local container_name=$1
	while true; do
		docker compose logs $container_name | grep "All components started successfully"
		if [ $? -eq 0 ]; then
			logt "Operator ${container_name} is ready for accept new task, proceeding to the next step..."
			break
		fi
		logt "Operator ${container_name} is not ready for accept new task, retrying in 2 second..."
		sleep 2
	done
	## Wait for operator to be ready
	sleep 3
}

check_operator_node_ready operator01
check_operator_node_ready operator02
echo -e "\n\n";

docker compose run --rm test;
STATUS=$?;
if [ "$STATUS" -ne 0 ]; then
	docker compose logs dvs -n 30;
	echo  -e "\n\n";
	echo  -e "\n\n";
	docker compose logs task-gateway -n 30;
	echo  -e "\n\n";
	echo  -e "\n\n";
	docker compose logs operator01 -n 30;
	echo  -e "\n\n";
	echo  -e "\n\n";
	docker compose logs operator02 -n 30;
	echo  -e "\n\n";
	echo  -e "\n\n";
	logt "Test failed";
	exit $STATUS;
fi

echo -e "\n\n"
echo -e "\n\n"
logt "Test passed for check reponse from chain";
logt "Start to check all operators task response";
echo -e "\n\n"

./scripts/check-task-response.sh;
STATUS_CHECK_ALL_OPERATORS=$?;
if [ "$STATUS_CHECK_ALL_OPERATORS" -ne "0" ]; then
		docker compose logs operator01 -n 20;
		echo  -e "\n\n";
		echo  -e "\n\n";
		docker compose logs operator02 -n 20;
		echo  -e "\n\n";
		echo  -e "\n\n";
		logt "Task response check failed";
		exit $STATUS_CHECK_ALL_OPERATORS;
fi

echo -e "\n\n"
logt "Test passed";
echo -e "\n\n"
exit 0

