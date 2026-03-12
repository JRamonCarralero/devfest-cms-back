package repository

import (
	"context"
	"devfest/internal/infrastructure/storage/dbgen"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	testPool    *pgxpool.Pool
	testQueries *dbgen.Queries
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	// 1. Init Postgres Container
	dbName := "testdb"
	dbUser := "user"
	dbPass := "pass"

	postgresContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPass),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start container: %s", err)
	}

	// 2. Obtain dinamic connection string
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %s", err)
	}

	// 3. Connect using pool
	testPool, err = pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to test db: %s", err)
	}

	// 4. Execute migrations
	if err := runMigrations(testPool); err != nil {
		log.Fatalf("failed to run migrations: %s", err)
	}

	// 5. Initialize queries from sqlc
	testQueries = dbgen.New(testPool)

	// 6. Excecute tests
	code := m.Run()

	// 7. Cleanup
	testPool.Close()
	postgresContainer.Terminate(ctx)

	os.Exit(code)
}

func runMigrations(db *pgxpool.Pool) error {
	path := "../migrations"
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, filepath.Join(path, f.Name()))
		}
	}
	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if _, err := db.Exec(context.Background(), string(content)); err != nil {
			return fmt.Errorf("error in %s: %w", file, err)
		}
	}
	return nil
}
