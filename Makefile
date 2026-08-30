include .env

run:
	@go run cmd/api/env/dev/main/main.go

build:
	@echo ==========================================
	@echo Building executable for Windows...
	@echo Reason: Windows Application Control blocks
	@echo executables generated in temporary folders
	@echo (like go run uses in AppData\Local\go-build).
	@echo So we build the binary locally to avoid that.
	@echo ==========================================
	go build -o app.exe cmd/api/env/dev/main/main.go
	.\app.exe


# Run all tests recursively.
# Only `go test ./...` is required.
# `findstr` is used on Windows to hide "[no test files]" lines.
test:
	go test ./... | findstr /v "[no test files]"

compose-dev:
	@docker-compose -f docker-compose.dev.yml --env-file .env up


## Other optional commands to improve the development experience 
git-diff:
	@git diff --staged -- . ':(exclude)go.mod' ':(exclude)go.sum' ':(exclude)*_templ.go' ':(exclude)web/web-app/resources/static/css/app.css'