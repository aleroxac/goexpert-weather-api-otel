conf ?= .env
include $(conf)
export $(shell sed 's/=.*//' $(conf))



## ---------- UTILS
.PHONY: help
help: ## Show this menu
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: clean
clean: ## Clean all temp files
	@rm -f coverage.*
	@docker image rm -f goexpert-weather-api-otel-input-api:latest
	@docker image rm -f goexpert-weather-api-otel-orchestrator-api:latest



## ----- COMPOSE
.PHONY: up
up: ## Put the compose containers up
	@docker-compose up -d

.PHONY: down
down: ## Put the compose containers down
	@docker-compose down



## ---------- MAIN
.PHONY: run
run: ## Make a request to the API
	@echo -e -----------------" input-api -----------------"
	@echo -n "422: "; curl -s "http://localhost:8080/cep" -d '{"cep": "1234567"}'
	@echo -n "200: "; curl -s "http://localhost:8080/cep" -d '{"cep": "13330250"}'

	@echo -e "\n\n------------- orchestrator-api -------------"
	@echo -n "422: "; curl -s "http://localhost:8081/cep/0100100"
	@echo -n "404: "; curl -s "http://localhost:8081/cep/01001009"
	@echo -n "200: "; curl -s "http://localhost:8081/cep/01001001"

.PHONY: test
test: ## Run the tests
	@go test -v ./... -coverprofile=coverage.out
	@go tool cover -html coverage.out -o coverage.html
