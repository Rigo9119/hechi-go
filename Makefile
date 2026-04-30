.PHONY: generate build run tidy db-up db-down db-logs migrate-up migrate-down

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -o bin/hechi-go .

run:
	go run .

tidy:
	go mod tidy

db-up:
	docker compose up -d db

db-down:
	docker compose down

db-logs:
	docker compose logs -f db

migrate-up: db-up
	@echo "Waiting for database to be ready..."
	@docker compose exec db sh -c 'until pg_isready -U hechi -d hechi; do sleep 1; done'
	psql "postgres://hechi:hechi@localhost:5432/hechi?sslmode=disable" -f migrations/001_initial.up.sql

migrate-down:
	psql "postgres://hechi:hechi@localhost:5432/hechi?sslmode=disable" -f migrations/001_initial.down.sql
