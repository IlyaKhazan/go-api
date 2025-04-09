include .env
export $(shell sed 's/=.*//' .env)

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	goose -dir database/migrations postgres "$(DB_URL)" up

migrate-add:
	goose -dir database/migrations create $(name) sql

migrate-down:
	goose -dir database/migrations postgres "$(DB_URL)" down

migrate-status:
	goose -dir database/migrations postgres "$(DB_URL)" status

.PHONY: migrate-up migrate-add migrate-down migrate-status
