package api

import (
	"devfest/internal/infrastructure/api/handlers"
	"devfest/internal/infrastructure/storage/dbgen"
	"devfest/internal/infrastructure/storage/repository"
	"devfest/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "database": "Connected"})
	})

	queries := dbgen.New(dbPool)

	// --- Dependencies Injection ---

	// Events
	eventRepo := repository.NewPostgresEventRepository(queries)
	eventUsecase := usecase.NewEventInteractor(eventRepo)
	eventHandler := handlers.NewEventHandler(eventUsecase)

	// ---ROUTES ---

	api := r.Group("/api/v1")
	{
		events := api.Group("/events")
		{
			events.GET("", eventHandler.GetEvents)       // Listar y buscar
			events.GET("/:slug", eventHandler.GetBySlug) // Ver detalle
		}
	}

	return r
}
