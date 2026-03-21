-include .env

# Variables
BINARY_NAME=server
MAIN_PATH=cmd/server/main.go
SCHEMA_PATH=internal/infrastructure/storage/migrations/000001_init_schema.up.sql

DB_URL=postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

.PHONY: generate tidy dev build db-init guard-env

# Validate .env
guard-env:
	@if [ ! -f .env ]; then \
		echo "FATAL: .env file not found, please copy .env.example"; \
		exit 1; \
	fi

# Initialize database (clear and recreate)
db-init: guard-env
	@echo "Recreating database $(DB_NAME) in $(DB_HOST)..."
	@psql $(DB_URL) -f $(SCHEMA_PATH)
	@echo "Database $(DB_NAME) ready!"

# Generate SQLC code
generate:
	@echo "Generating database code with SQLC..."
	./sqlc.exe generate

# Download dependencies
tidy:
	go mod tidy

# Run server in development mode
dev: generate
	go run $(MAIN_PATH)

# Build server
build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

# Prepare environment
prep: db-init generate
	@echo "🚀 Database updates and Go code generated."