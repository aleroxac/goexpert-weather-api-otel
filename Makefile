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



## ---------- MAIN
.PHONY: up
up: ## put the docker-compose containers up
	@docker-compose up -d

.PHONY: down
down: ## put the docker-compose containers down
	@docker-compose down

.PHONY: run
run: ## make a request to the API
	@curl -sv -X POST http://localhost:8080/weather -H "Content-Type: application/json" -d '{"cep": "13330250"}'
