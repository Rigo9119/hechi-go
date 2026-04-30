.PHONY: generate build run migrate-up migrate-down tidy

generate:
	go run github.com/99designs/gqlgen generate

build: generate
	go build -o bin/hechi-go .

run:
	go run .

tidy:
	go mod tidy

migrate-up:
	psql "$(DATABASE_URL)" -f migrations/001_initial.up.sql

migrate-down:
	psql "$(DATABASE_URL)" -f migrations/001_initial.down.sql
