.PHONY: help run run-dev build migrate-up migrate-down migrate-create test clean init-docker

DATABASE_URL ?= postgres://postgres:postgres@localhost:54323/practiq-db?sslmode=disable

help:
	@echo "Available commands:"
	@echo "  make run              - Run the API server"
	@echo "  make run-dev          - Run with Air live reloader"
	@echo "  make build            - Build the binary"
	@echo "  make migrate-up       - Apply migrations"
	@echo "  make migrate-down     - Revert last migration"
	@echo "  make migrate-create   - Create a new migration"
	@echo "  make test             - Run tests"
	@echo "  make init-docker      - Start Docker services"

run:
	go run ./cmd/api/

run-dev:
	air

build:
	go build -o ./build/practiq-be ./cmd/api/

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	@read -p "Migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

test:
	go test -coverprofile=coverage.out ./...

clean:
	rm -rf build/ coverage.out tmp/

init-docker:
	docker-compose up -d

install-deps:
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
