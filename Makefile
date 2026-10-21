infra-up:
	@docker-compose up -d

ptoto-generate:
	@protoc -I proto proto/order/.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go --go-grpc_opt=paths=source_relative 

infra-logs: ## Посмотреть логи общей инфраструктуры
	@docker compose logs -f    
