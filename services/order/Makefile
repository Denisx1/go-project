ifeq ($(OS),Windows_NT)
    PROJECT_ROOT = $(CURDIR)
	RM_PGDATA = if exist out\pgdata rmdir /s /q out\pgdata
	RM_LOGS   = if exist out\log rmdir /s /q out\log
	 
else
    PROJECT_ROOT = $(shell pwd)
	RM_PGDATA = rm -rf ./out/pgdata
	RM_LOGS = rm -rf ./out/log
endif

include .env
CONTRACTS_DIR := ../../contracts
ENV_FILE:=.env
export

env-up:
	@docker-compose --env-file $(ENV_FILE) up -d $(POSTGRES_CONTAINER)
env-down:
	@docker-compose --env-file $(ENV_FILE) down $(POSTGRES_CONTAINER)
env-cleanup:
	@docker compose --env-file $(ENV_FILE) down -v $(POSTGRES_CONTAINER)
	@$(RM_PGDATA)
	@echo "All data cleaned up for $(SERVICE_NAME)."

env-port-forward: 
	@docker compose --env-file $(ENV_FILE)  up -d $(PORT_FORWARDER_CONTAINER)

env-port-close:
	@docker compose --env-file $(ENV_FILE)  stop $(PORT_FORWARDER_CONTAINER)
	@docker compose --env-file $(ENV_FILE)  rm -f $(PORT_FORWARDER_CONTAINER)

migrate-create:
ifndef seq
	$(error ❌ Ошибка: Не указано имя миграции! Использование: make migrate-create s=$(s) seq=add_users_table)
endif
	@docker-compose --env-file $(ENV_FILE) run --rm $(MIGRATION_CONTAINER) \
	create \
	-ext sql \
	-dir ${MIGRATION_FOLDER} \
	-seq "$(seq)"

logs-cleanup:
	@$(RM_LOGS)
	@echo "All data cleaned up for $(SERVICE_NAME)."
migrate-up:
	@make migrate-action action="$(if $(action),$(action),up)"
migrate-down:
	@make migrate-action action="$(if $(action),$(action),down 1)"
migrate-drop:
	@make migrate-action action="drop -f"
migrate-version:
	@make migrate-action action=version
proto-generate:
	@protoc -I $(CONTRACTS_DIR)/proto \
		$(CONTRACTS_DIR)/proto/$(SERVICE_NAME)/$(SERVICE_NAME).proto \
		--go_out=$(CONTRACTS_DIR)/gen/go \
		--go_opt=paths=source_relative \
		--go-grpc_out=$(CONTRACTS_DIR)/gen/go \
		--go-grpc_opt=paths=source_relative
migrate-action:
	@docker-compose --env-file $(ENV_FILE) run --rm $(MIGRATION_CONTAINER) \
		-path /migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_CONTAINER):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		$(action)

run: ## Запустить Go-код сервиса локально
	@go run ./cmd/order