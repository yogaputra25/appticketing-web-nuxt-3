.PHONY: help build run test migrate-up migrate-down seed smoke tidy

DATABASE_URL ?= postgres://ticketing:ticketing_secret@localhost:5432/ticketing?sslmode=disable

help:
	@echo "Targets:"
	@echo "  build         - compile binary"
	@echo "  run           - run server"
	@echo "  test          - run unit tests"
	@echo "  migrate-up    - apply migrations"
	@echo "  migrate-down  - rollback last migration"
	@echo "  seed          - run seeder"
	@echo "  smoke         - run smoke test script (requires running API)"
	@echo "  tidy          - go mod tidy"

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

seed:
	go run ./cmd/seed

smoke:
	@echo "Running smoke test (API must be running)..."
	@bash scripts/smoke.sh

tidy:
	go mod tidy
