NET_NAME := shared-network

init-net:
	@docker network inspect $(NET_NAME) >/dev/null 2>&1 || docker network create $(NET_NAME) >/dev/null

infra-up: init-net
	@docker compose up -d

infra-down:
	@docker compose down
infra-logs: ## Посмотреть логи общей инфраструктуры
	@docker compose logs -f    	

.PHONY: init-net infra-up infra-down