package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresClient(dsn string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("error parsing DSN: %w", err)
	}

	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	config.MaxConns = 10
	config.MinConns = 2

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("connection to PostgreSQL (Supabase) failed creating pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping to PostgreSQL (Supabase) failed: %w", err)
	}

	log.Println("✅ Connection to PostgreSQL (Supabase) established successfully")
	return pool, nil
}

func InitializeSchema(db *pgxpool.Pool, sqlFilePath string) error {
	log.Printf("Reading SQL schema from %s...", sqlFilePath)

	content, err := os.ReadFile(sqlFilePath)
	if err != nil {
		return fmt.Errorf("could not read sql file: %w", err)
	}

	_, err = db.Exec(context.Background(), string(content))
	if err != nil {
		return fmt.Errorf("could not execute schema: %w", err)
	}

	log.Println("✅ Database schema initialized/verified successfully")
	return nil
}

func BuildDSN(user, pass, host, port, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&pgbouncer=true",
		url.PathEscape(user),
		url.PathEscape(pass),
		host,
		port,
		dbname,
	)
}
