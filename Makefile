.DEFAULT_GOAL := help

build: ## Build a binary
	mkdir -p build && CGO_ENABLED=1 go build -o build/main

.PHONY: test
test: ## Run automated tests
	go test ./...

.PHONY: lint ## Format code and tidy go modules
lint:
	go fmt
	go mod tidy

.PHONY: clean ## Clean temporary files
clean:

.PHONY: run ## Run main server
run:
	go run main.go server

MIGRATION_NAME?=migration
TS=$(shell date +%Y%m%d-%H%M%S)
.PHONY: add-migration
add-migration:
	@$(shell mkdir -p "src/db/ddl/migrations/${TS}-${MIGRATION_NAME}")
	@$(shell touch src/db/ddl/migrations/${TS}-${MIGRATION_NAME}/{up.sql,down.sql})

.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'