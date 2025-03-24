default:
	@echo "hello world"

check-env-gh-token:
	@if [ -z "$${GITHUB_TOKEN}" ] && ! grep -q '^GITHUB_TOKEN=' docker/.env 2>/dev/null; then \
		echo "Error: GITHUB_TOKEN is not set in environment or docker/.env file"; \
		exit 1; \
	else \
		echo "GITHUB_TOKEN is set."; \
	fi

docker-build-all: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build
	@echo "docker build all done, `date`"

docker-build-contracts: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build hardhat
	@echo "docker build contracts done, `date`"

docker-build-pelldvs: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build pelldvs
	@echo "docker build pelldvs done, `date`"

docker-build-operator: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build operator
	@echo "docker build operator done, `date`"

docker-all-up:
	@cd docker && docker compose up -d

docker-all-down:
	@cd docker && docker compose down -v

docker-all-status:
	@cd docker && docker compose ps -a

docker-up-operator:
	@cd docker && docker compose up operator01 operator02 -d

docker-mix-operators-up:
	@cd docker && docker compose up operator01 operator02 -d

docker-mix-operators-down:
	@cd docker && docker compose down operator01 operator02

docker-hardhat-up:
	@cd docker && docker compose up hardhat -d

docker-hardhat-down:
	@cd docker && docker compose down hardhat -v

docker-hardhat-logs:
	@cd docker && docker compose logs hardhat -f

docker-hardhat-shell:
	@cd docker && docker compose exec -it hardhat bash

docker-hardhat-rerun:
	make docker-hardhat-down
	make docker-hardhat-up
	make docker-hardhat-logs

docker-emulator-up:
	@cd docker && docker compose up emulator -d

docker-emulator-down:
	@cd docker && docker compose down emulator -v

docker-emulator-logs:
	@cd docker && docker compose logs -f emulator

docker-emulator-shell:
	@cd docker && docker compose exec -it emulator bash

docker-emulator-rerun:
	make docker-emulator-down
	make docker-emulator-up
	make docker-emulator-logs

docker-dvs-up:
	@cd docker && docker compose up dvs -d

docker-dvs-down:
	@cd docker && docker compose down dvs -v

docker-dvs-logs:
	@cd docker && docker compose logs dvs -f

docker-dvs-shell:
	@cd docker && docker compose exec -it dvs bash

docker-dvs-rerun:
	make docker-dvs-down
	make docker-dvs-up
	make docker-dvs-logs

docker-gateway-up:
	@cd docker && docker compose up task-gateway -d

docker-gateway-down:
	@cd docker && docker compose down task-gateway -v

docker-gateway-logs:
	@cd docker && docker compose logs task-gateway -f

docker-gateway-shell:
	@cd docker && docker compose exec -it task-gateway bash

docker-gateway-rerun:
	make docker-gateway-down
	make docker-gateway-up
	make docker-gateway-logs

docker-operator-all-up:
	@cd docker && docker compose up operator01 operator02 -d

docker-operator-all-down:
	@cd docker && docker compose down operator01 operator02 -v

docker-operator-all-logs:
	@cd docker && docker compose logs operator01 operator02 -f

docker-operator-shell:
	@cd docker && docker compose exec -it operator bash

docker-operator-all-rerun:
	make docker-operator-all-down
	make docker-operator-all-up
	make docker-operator-all-logs

docker-operator-01-up:
	@cd docker && docker compose up operator01  -d

docker-operator-01-down:
	@cd docker && docker compose down operator01  -v

docker-operator-01-logs:
	@cd docker && docker compose logs operator01  -f

docker-operator-01-shell:
	@cd docker && docker compose exec -it operator01 bash

docker-operator-01-rerun:
	make docker-operator-01-down
	make docker-operator-01-up
	make docker-operator-01-logs

docker-test:
	@bash ./docker/scripts/test-in-host.sh

test:
	@go test ./...

build:
	@go build -mod=readonly -o bin/squaringd ./cmd/squaringd

#? lint: Run latest golangci-lint linter
lint:
	@echo "--> Running linter"
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run
.PHONY: lint

# Run goimports-reviser to lint and format imports
lint-imports:
	@find . -name "*.go" -not -path "./vendor/*" -not -path "./.git/*" | while read -r file; do \
		goimports-reviser -rm-unused -format "$$file"; \
	done
	
.PHONY: proto
proto:
	@rm -rf dvs/query/types/msg_service.pb.gw.go
	@rm -rf dvs/query/msg_service.pb.gw.go
	@rm -rf dvs/query/types/*.pb.go
	@rm -rf dvs/types/*.pb.go
	@cd proto && buf generate --template buf.gen.gogo.yaml
	@mv -f dvs/query/msg_service.pb.gw.go dvs/query/types/