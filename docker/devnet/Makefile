default:
	@echo "hello"

check-env-gh-token:
	@if [ -z "$${GITHUB_TOKEN}" ] && ! grep -q '^GITHUB_TOKEN=' ../.env 2>/dev/null; then \
		echo "Error: GITHUB_TOKEN is not set in environment or ../.env file"; \
		exit 1; \
	else \
		echo "GITHUB_TOKEN is set."; \
	fi

docker-testnet-build-all: check-env-gh-token
	docker compose -p pelldvs-example-testnet -f ../docker-compose.build.yml build

docker-testnet-build-pelldvs: check-env-gh-token
	docker compose -p pelldvs-example-testnet -f ../docker-compose.build.yml build pelldvs

docker-testnet-build-operator: check-env-gh-token
	docker compose -p pelldvs-example-testnet -f ../docker-compose.build.yml build operator

dvs-up-all:
	make dvs-up-dvs;
	make dvs-up-gateway;

dvs-up-dvs:
	docker compose --env-file .env-dvs -f docker-compose.dvs.yml up dvs -d

dvs-up-gateway:
	docker compose --env-file .env-gateway -f docker-compose.dvs.yml up gateway -d

dvs-logs-dvs:
	docker compose -f docker-compose.dvs.yml logs -f dvs -n 100

dvs-logs-gateway:
	docker compose -f docker-compose.dvs.yml logs -f gateway -n 100

dvs-down-all:
	@make dvs-down-dvs;
	@make dvs-down-gateway;

dvs-down-dvs:
	docker compose --env-file .env-dvs -f docker-compose.dvs.yml down -v dvs

dvs-down-gateway:
	docker compose --env-file .env-gateway -f docker-compose.dvs.yml down -v gateway

dvs-rm-dvs:
	docker compose --env-file .env-dvs -f docker-compose.dvs.yml rm dvs -f

dvs-rm-gateway:
	docker compose --env-file .env-gateway -f docker-compose.dvs.yml rm gateway -f

dvs-dvs-rerun:
	make dvs-down-dvs;
	make dvs-up-dvs;
	make dvs-logs-dvs;

dvs-gateway-rerun:
	make dvs-down-gateway;
	make dvs-up-gateway;
	make dvs-logs-gateway;

operator-up-operator01:
	docker compose --env-file .env-operator01 -f docker-compose.yml up operator01 -d

operator-up-operator02:
	docker compose --env-file .env-operator02 -f docker-compose.yml up operator02 -d

operator-up-all:
	@make operator-up-operator01;
	@make operator-up-operator02;

operator-down-operator01:
	docker compose --env-file .env-operator01 -f docker-compose.yml down -v

operator-down-operator02:
	docker compose --env-file .env-operator02 -f docker-compose.yml down -v

operator-down-all:
	@make operator-down-operator01;
	@make operator-down-operator02;

operator-logs-operator01:
	docker compose -f docker-compose.yml logs -f operator01 -n 100

operator-logs-operator02:
	docker compose -f docker-compose.yml logs -f operator02 -n 100

