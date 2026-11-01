include .env

CONTRACTS_DIR := ../../contracts
ENV_FILE := .env
DC := docker compose --env-file $(ENV_FILE)
export

env-up:
	@$(DC) up -d $(POSTGRES_CONTAINER)

env-down:
	@$(DC) stop $(POSTGRES_CONTAINER)

env-cleanup:
	@$(DC) down -v
	@rm -rf ./out/pgdata
	@echo "All data cleaned up for $(SERVICE_NAME)."

env-port-forward:
	@$(DC) up -d $(POSTGRES_PORT_FORWARDER)

env-port-close:
	@$(DC) stop $(POSTGRES_PORT_FORWARDER)
	@$(DC) rm -f $(POSTGRES_PORT_FORWARDER)

migrate-create:
ifndef seq
	$(error Не указано имя миграции! Использование: make migrate-create seq=add_users_table)
endif
	@$(DC) run --rm $(MIGRATION_CONTAINER) \
		create \
		-ext sql \
		-dir $(MIGRATION_FOLDER) \
		-seq "$(seq)"

logs-cleanup:
	@rm -rf ./out/log
	@echo "Logs cleaned up for $(SERVICE_NAME)."

migrate-up:
	@$(MAKE) migrate-action action="$(if $(action),$(action),up)"

migrate-down:
	@$(MAKE) migrate-action action="$(if $(action),$(action),down 1)"

migrate-drop:
	@$(MAKE) migrate-action action="drop -f"

migrate-version:
	@$(MAKE) migrate-action action=version

migrate-action:
	@$(DC) run --rm $(MIGRATION_CONTAINER) \
		-path $(MIGRATION_FOLDER) \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_CONTAINER):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		$(action)
		
proto-generate:
	@$(MAKE) -C ../../contracts proto-generate

run:
	@go run ./cmd/order
