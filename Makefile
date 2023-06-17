PROJECT_NAME:=altt
FILE_HASH := $(shell git rev-parse HEAD)
GOLANGCI_LINT := $(shell command -v golangci-lint 2> /dev/null)

init_repo: ## create necessary configs
	cp configs/sample.app_conf.yml configs/app_conf.yml
	cp configs/sample.app_conf_docker.yml configs/app_conf_docker.yml

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-lint: ## Installs golangci-lint tool which a go linter
ifndef GOLANGCI_LINT
	${info golangci-lint not found, installing golangci-lint@latest}
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
endif

abi: ## generate abi struct
	abigen --abi internal/service/web3/approver/erc20.abi.json --pkg approver --type Erc20 --out erc_20.go
	mv erc_20.go internal/service/web3/approver/
	abigen --abi internal/service/web3/swapper/stargate.abi.json --pkg swapper --type StargateRouter --out stargate_abi.go
	mv stargate_abi.go internal/service/web3/swapper/

gogen: ## generate code
	${info generate code...}
	go generate ./internal...

test: ## Runs tests
	${info Running tests...}
	go test -v -race ./... -cover -coverprofile cover.out
	go tool cover -func cover.out | grep total

bench: ## Runs benchmarks
	${info Running benchmarks...}
	go test -bench=. -benchmem ./... -run=^#

vulcheck: ## Runs vulnerability check
	${info Running vulnerability check...}
	govulncheck ./...

lint: install-lint ## Runs linters
	@echo "-- linter running"
	golangci-lint run -c .golangci.yaml ./internal...
	golangci-lint run -c .golangci.yaml ./cmd...

stop: ## Stops the local environment
	${info Stopping containers...}
	docker container ls -q --filter name=${PROJECT_NAME} ; true
	${info Dropping containers...}
	docker rm -f -v $(shell docker container ls -q --filter name=${PROJECT_NAME}) ; true

dev_up: init_repo stop ## Runs local environment
	${info Running docker-compose up...}
	GIT_HASH=${FILE_HASH} docker compose -p ${PROJECT_NAME} up --build

build: ## Builds binary
	@echo "-- building binary"
	go build -o ./bin/binary ./cmd


.PHONY: help install-lint test gogen lint stop dev_up build init_repo vulcheck
.DEFAULT_GOAL := help
