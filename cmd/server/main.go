package main

import (
	db "devfest/internal/infrastructure/storage"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database Connection

	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require&pgbouncer=true",
		url.PathEscape(user),
		url.PathEscape(pass),
		host,
		dbPort,
		dbname,
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

	// ToDo: implement routes and inject dependencies

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "database": "Connected"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on port: %s", port)
	r.Run(":" + port)
}
