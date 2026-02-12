include .env

run:
	@go run cmd/api/env/dev/main.go

build:
	@echo ==========================================
	@echo Building executable for Windows...
	@echo Reason: Windows Application Control blocks
	@echo executables generated in temporary folders
	@echo (like go run uses in AppData\Local\go-build).
	@echo So we build the binary locally to avoid that.
	@echo ==========================================
	go build -o app.exe cmd/api/env/dev/main.go
	.\app.exe

compose-dev:
	@docker-compose -f docker-compose.dev.yml --env-file .env up
