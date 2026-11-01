include .env

CONTRACTS_DIR := ../../contracts
ENV_FILE := .env
DC := docker compose --env-file $(ENV_FILE)
export

env-up:
	@$(DC) up -d $(POSTGRES_CONTAINER)
	@$(DC) up -d $(REDIS_CONTAINER)

env-down:
	@$(DC) stop $(POSTGRES_CONTAINER)
	@$(DC) stop $(REDIS_CONTAINER)

env-cleanup:
	@$(DC) down -v
	@rm -rf ./out/pgdata
	@rm -rf ./out/redisdata
	@echo "All data cleaned up for $(SERVICE_NAME)."

env-port-forward:
	@$(DC) up -d $(POSTGRES_PORT_FORWARDER)
	@$(DC) up -d $(REDIS_PORT_FORWARDER)
	

env-port-close:
	@$(DC) stop $(POSTGRES_PORT_FORWARDER)
	@$(DC) stop $(REDIS_PORT_FORWARDER)
	@$(DC) rm -f $(POSTGRES_PORT_FORWARDER)
	@$(DC) rm -f $(REDIS_PORT_FORWARDER)

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
migrate-force:
	@$(MAKE) migrate-action action="force $(if $(step),$(step),1)"
migrate-version:
	@$(MAKE) migrate-action action=version

migrate-action:
	@export UID=$$(id -u); export GID=$$(id -g); 
	@$(DC) run --rm $(MIGRATION_CONTAINER) \
		-path $(MIGRATION_FOLDER) \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_CONTAINER):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		$(action)
		
proto-generate:
	@$(MAKE) -C ../../contracts proto-generate

run:
	@go run ./cmd/user
