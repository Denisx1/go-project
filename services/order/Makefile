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
	@$(DC) up -d $(PORT_FORWARDER_CONTAINER)

env-port-close:
	@$(DC) stop $(PORT_FORWARDER_CONTAINER)
	@$(DC) rm -f $(PORT_FORWARDER_CONTAINER)

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

.PHONY: env-up env-down env-cleanup env-port-forward env-port-close \
	migrate-create logs-cleanup migrate-up migrate-down migrate-drop \
	migrate-version migrate-action proto-generate run