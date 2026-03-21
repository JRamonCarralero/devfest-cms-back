-include .env

# Variables
BINARY_NAME=server
MAIN_PATH=cmd/server/main.go
SCHEMA_PATH=internal/infrastructure/storage/migrations/000001_init_schema.up.sql

.PHONY: generate tidy dev build db-init guard-env

# Validate .env
guard-env:
	@if [ ! -f .env ]; then \
		echo "FATAL: .env file not found, please copy .env.example"; \
		exit 1; \
	fi

# Initialize database (clear and recreate)
db-init: guard-env
	@echo "Configuring database $(DB_USER)..."
	@echo $(DB_HOST):$(DB_PORT_MF):$(DB_NAME):$(DB_USER):$(DB_PASSWORD) > pgpass_temp
	@echo "Conectando a Supabase..."
	@set "PGPASSFILE=pgpass_temp" && psql -h $(DB_HOST) -p $(DB_PORT_MF) -U $(DB_USER) -d $(DB_NAME) -f $(SCHEMA_PATH)
	@rm -f pgpass_temp || del pgpass_temp
	@echo "✅ Tables initialized!"

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