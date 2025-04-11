name ?= init_schema

# Создание новой миграции
migrate-add: ## Create new migration file, usage: make migrate-add name=<migration_name>
	@echo "Creating migration: $(name)"
	goose create --dir migrations $(name) sql

# Применить миграции (через твой migrate.go)
migrate-up: ## Apply all migrations via Go code
	go run ./migrations

# (Если хочешь позже — можно добавить migrate-down тоже через Go)
