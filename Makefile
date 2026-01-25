include .env

run:
	@go run cmd/api/env/dev/main.go

compose-dev:
	@docker-compose -f docker-compose.dev.yml --env-file .env up