package main

import (
	"devfest/internal/infrastructure/api"
	"devfest/internal/infrastructure/config"
	db "devfest/internal/infrastructure/storage"
	"log"
)

func main() {
	// Load Environment Variables
	cfg := config.LoadConfig()

	// Database Connection
	dsn := db.BuildDSN(
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	dbPool, err := db.NewPostgresClient(dsn)
	if err != nil {
		log.Fatalf("❌ Error: %v", err)
	}
	defer dbPool.Close()

	log.Println("🚀 Connection secure and success!")

	// Initialize Schema

	sqlPath := "internal/infrastructure/storage/migrations/000001_init_schema.up.sql"

	if err := db.InitializeSchema(dbPool, sqlPath); err != nil {
		log.Fatalf("❌ Error: %v", err)
	}

	r := api.SetupRouter(dbPool)

	log.Printf("🚀 Server running on port: %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
