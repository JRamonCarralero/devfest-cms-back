package main

import (
	"devfest/internal/infrastructure/api"
	db "devfest/internal/infrastructure/storage"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database Connection

	dsn := db.BuildDSN(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port: %s", port)
	r.Run(":" + port)
}
