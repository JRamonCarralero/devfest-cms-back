# Community Event API (Go + Clean Architecture)

A robust and scalable backend service built with **Go** to manage technology community events. This project follows **Clean Architecture** principles and leverages **Supabase (PostgreSQL)** for reliable data persistence.

## 🚀 Tech Stack

- **Language:** Go 1.25+
- **Database:** PostgreSQL (via Supabase)
- **SQL Generator:** [SQLC](https://sqlc.dev/) (Type-safe SQL)
- **Driver:** [pgx/v5](https://github.com/jackc/pgx) (PostgreSQL Driver and Toolkit)
- **Architecture:** Clean Architecture (Entities, Repository Pattern, Dependency Injection)

## 📁 Project Structure

```text
├── cmd/
│   └── server/             # Application entry point (main.go)
├── internal/
│   ├── domain/             # Core business logic & Interfaces (Entities)
│   ├── infrastructure/
│   │   ├── storage/        # Database implementations & migrations
│   │   │   ├── dbgen/      # SQLC GENERATED CODE (Do not edit manually)
│   │   │   ├── queries/    # Pure SQL queries (.sql files)
│   │   │   ├── migrations/ # Table definitions & DB triggers
│   │   │   └── repository/ # Repository Pattern & Dependency Injection
│   │   └── config/         # Environment variable management
└── sqlc.yaml               # SQLC configuration for Go code generation
```

## 🛠️ Getting Started

### Prerequisites

- Go SDK installed.

- SQLC binary (Download sqlc.exe for Windows or install via brew/scoop).

A running Supabase or PostgreSQL instance.

###  1. Environment Configuration

Create a .env file in the root directory and provide your database credentials:

Fragmento de código
DB_USER=your_user
DB_PASSWORD=your_password
DB_HOST=your_host
DB_PORT=5432
DB_NAME=postgres

### 2. Database Code Generation

This project uses SQLC to generate type-safe Go code from raw SQL. If you modify any .sql file inside internal/infrastructure/storage/queries/, you must regenerate the code:

```Bash
# Using Makefile
make generate

# Or manually using the binary
./sqlc.exe generate
```

### 3. Running the Server

To start the development server with automatic generation:

```Bash
make dev
```

## 📝 Key Features

- **Type Safety**: SQLC ensures that all database operations are checked at compile time.

- **Automated Auditing**: Primary tables include database-level triggers to manage updated_at timestamps automatically.

- **Clean Separation**: The business logic (Domain) is completely decoupled from the database implementation (Infrastructure).

- **Pagination**: Built-in support for paginated results and total count metadata for frontend integration (e.g., Astro).

## 🤝 Contribution

1. Ensure your SQL queries are formatted correctly.

2. Run make tidy to clean up dependencies before submitting a PR.

3. Documentation for new endpoints should be added to the relevant domain files.
