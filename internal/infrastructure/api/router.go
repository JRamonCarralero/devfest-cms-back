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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(dbPool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	r.Use(middleware.TraceMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP", "database": "Connected"})
	})

	r.StaticFile("/docs/swagger.yaml", "./api-docs/swagger.yaml")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/docs/swagger.yaml"),
	))

	queries := dbgen.New(dbPool)

	api := r.Group("/api/v1")

	// --- Dependencies Injection ---

	// Events
	eventRepo := repository.NewPostgresEventRepository(queries)
	eventUsecase := usecase.NewEventInteractor(eventRepo)
	eventHandler := handlers.NewEventHandler(eventUsecase)

	// Persons
	personRepo := repository.NewPostgresPersonRepository(queries)
	personUsecase := usecase.NewPersonInteractor(personRepo)
	personHandler := handlers.NewPersonHandler(personUsecase)

	// ---ROUTES ---

	{
		events := api.Group("/events")
		protecteEvents := events.Group("/")
		protecteEvents.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
		{
			events.GET("", eventHandler.GetEvents)               // All events
			events.GET("/id/:id", eventHandler.GetByID)          // Event by ID
			events.GET("/slug/:slug", eventHandler.GetBySlug)    // Event by slug
			events.GET("/status/active", eventHandler.GetActive) // All active events
			events.GET("/paged", eventHandler.GetPaged)          // All events paged

			protecteEvents.POST("", eventHandler.Create)       // Create event
			protecteEvents.PUT("/:id", eventHandler.Update)    // Update event
			protecteEvents.DELETE("/:id", eventHandler.Delete) // Delete event
		}
	}

	persons := api.Group("/persons")
	protectedPersons := persons.Group("/")
	protectedPersons.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		persons.GET("", personHandler.GetAll)
		persons.GET("/paged", personHandler.GetPaged)
		persons.GET("/id/:id", personHandler.GetByID)
		persons.GET("/email/:email", personHandler.GetByEmail)

		protectedPersons.POST("", personHandler.Create)
		protectedPersons.PUT("/:id", personHandler.Update)
		protectedPersons.DELETE("/:id", personHandler.Delete)
	}

	return r
}
