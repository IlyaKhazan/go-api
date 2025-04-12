-include .env
export $(shell sed 's/=.*//' .env 2>/dev/null)

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

migrate-up:
	docker-compose exec go-gin-api goose -dir database/migrations postgres "$(DB_URL)" up

migrate-down:
	docker-compose exec go-gin-api goose -dir database/migrations postgres "$(DB_URL)" down

migrate-status:
	docker-compose exec go-gin-api goose -dir database/migrations postgres "$(DB_URL)" status

migrate-add:
	docker-compose exec go-gin-api goose -dir database/migrations create create_users_table sql

.PHONY: migrate-up migrate-add migrate-down migrate-status
