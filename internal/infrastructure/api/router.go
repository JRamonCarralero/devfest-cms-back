package api

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/handlers"
	"devfest/internal/infrastructure/api/middleware"
	"devfest/internal/infrastructure/storage/dbgen"
	"devfest/internal/infrastructure/storage/repository"
	"devfest/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.TraceMiddleware())

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
		protecteEvents := events.Group("/")
		protecteEvents.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
		{
			events.GET("", eventHandler.GetEvents)            // All events
			events.GET("/:id", eventHandler.GetByID)          // Event by ID
			events.GET("/slug/:slug", eventHandler.GetBySlug) // Event by slug
			events.GET("/active", eventHandler.GetActive)     // All active events
			events.GET("/page", eventHandler.GetPaged)        // All events paged

			protecteEvents.POST("", eventHandler.Create)       // Create event
			protecteEvents.PUT("/:id", eventHandler.Update)    // Update event
			protecteEvents.DELETE("/:id", eventHandler.Delete) // Delete event
		}
	}

	return r
}
