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

docker-build-contracts: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build hardhat

docker-build-pelldvs: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build pelldvs

docker-build-operator: check-env-gh-token
	@cd docker && docker compose -f docker-compose.build.yml build operator

docker-all-up:
	@cd docker && docker compose up -d

docker-all-down:
	@cd docker && docker compose down -v

docker-all-status:
	@cd docker && docker compose ps -a

docker-up-operator:
	@cd docker && docker compose up operator01 operator02 -d

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
	@cd proto && buf generate --template buf.gen.gogo.yaml