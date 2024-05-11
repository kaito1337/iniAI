#!/usr/bin/make -f
.DEFAULT_GOAL := help

include .env

help:
	@echo " "
	@echo "iniVoice MVP"
	@echo " "
	@echo "Usage:"
	@echo " "
	@echo "> make run - Run iniVoice BOT"
	@echo " "
	@echo "> make dev - Run iniVoice BOT in dev mode"
	@echo " "
	@echo "> make db-status - Check database status"
	@echo " "
	@echo "> make db-up - Create database"
	@echo " "
	@echo "> make db-down - Drop database"
	@echo " "
	@echo "> make go-exports - Configure AnyCall File Handle Service"
	@echo " "
	@echo "Dependencies:"
	@echo " "
	@echo "Air (go install github.com/cosmtrek/air@latest)"
	@echo "Goose (go install github.com/pressly/goose/v3/cmd/goose@latest)"
	@echo " "
	@echo "Troubleshooting:"
	@echo " "
	@echo "Error: command not found: air || command not found: goose"
	@echo " "
	@echo "Solution:"
	@echo " - export GOPATH=#HOME/go"
	@echo " - export PATH=#GOPATH/bin:#PATH"
	@echo " -- change # -> $$"

run:
	@echo " > Running Backend Server"
	go run .\cmd\main\main.go


dev:
	@echo " > Running Backend Server in dev mode"
			air -c .air.toml

db-status:
	@echo " > Checking database status"
		goose -dir='./internal/db/migrations' postgres 'host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} sslmode=disable' status

db-up:
	@echo " > Creating database"
		goose -dir='./internal/db/migrations' postgres 'host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} sslmode=disable' up

db-down:
	@echo " > Dropping database"
		goose -dir='./internal/db/migrations' postgres 'host=${DB_HOST} port=${DB_PORT} user=${DB_USERNAME} password=${DB_PASSWORD} dbname=${DB_DATABASE} sslmode=disable' down-to -1

go-exports:
	@echo " > Configuring AnyCall Automation Resolver"
	bash -c "export GOPATH=$$HOME/go | export PATH=$$GOPATH/bin:$$PATH"
