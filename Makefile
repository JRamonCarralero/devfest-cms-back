# Variables
BINARY_NAME=server
MAIN_PATH=cmd/server/main.go

.PHONY: generate build run

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