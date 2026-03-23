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
	eventRepo := repository.NewEventRepository(queries)
	eventUsecase := usecase.NewEventInteractor(eventRepo)
	eventHandler := handlers.NewEventHandler(eventUsecase)

	// Persons
	personRepo := repository.NewPersonRepository(queries)
	personUsecase := usecase.NewPersonInteractor(personRepo)
	personHandler := handlers.NewPersonHandler(personUsecase)

	// Collaborators
	collaboratorRepo := repository.NewCollaboratorRepository(queries)
	collaboratorUsecase := usecase.NewCollaboratorInteractor(collaboratorRepo, personRepo, eventRepo)
	collaboratorHandler := handlers.NewCollaboratorHandler(collaboratorUsecase)

	// Developers
	developerRepo := repository.NewDeveloperRepository(queries)
	developerUsecase := usecase.NewDeveloperInteractor(developerRepo, personRepo, eventRepo)
	developerHandler := handlers.NewDeveloperHandler(developerUsecase)

	// Organizers
	organizerRepo := repository.NewOrganizerRepository(queries)
	organizerUsecase := usecase.NewOrganizerInteractor(organizerRepo, personRepo, eventRepo)
	organizerHandler := handlers.NewOrganizerHandler(organizerUsecase)

	// Speakers
	speakerRepo := repository.NewSpeakerRepository(queries)
	speakerUsecase := usecase.NewSpeakerInteractor(speakerRepo, personRepo, eventRepo)
	speakerHandler := handlers.NewSpeakerHandler(speakerUsecase)

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

	collaborators := api.Group("/collaborators")
	protectedCollaborators := collaborators.Group("/")
	protectedCollaborators.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		collaborators.GET("/event/:event-id", collaboratorHandler.GetAll)
		collaborators.GET("/id/:id", collaboratorHandler.GetByID)
		collaborators.GET("/event/:event-id/paged", collaboratorHandler.ListPaged)

		protectedCollaborators.POST("", collaboratorHandler.Create)
		protectedCollaborators.PUT("/:id", collaboratorHandler.Update)
		protectedCollaborators.DELETE("/:id", collaboratorHandler.Delete)
	}

	developers := api.Group("/developers")
	protectedDevelopers := developers.Group("/")
	protectedDevelopers.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		developers.GET("/event/:event-id", developerHandler.GetAll)
		developers.GET("/id/:id", developerHandler.GetByID)
		developers.GET("/event/:event-id/paged", developerHandler.ListPaged)

		protectedDevelopers.POST("", developerHandler.Create)
		protectedDevelopers.PUT("/:id", developerHandler.Update)
		protectedDevelopers.DELETE("/:id", developerHandler.Delete)
	}

	organizers := api.Group("/organizers")
	protectedOrganizers := organizers.Group("/")
	protectedOrganizers.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		organizers.GET("/event/:event-id", organizerHandler.GetAll)
		organizers.GET("/id/:id", organizerHandler.GetByID)
		organizers.GET("/event/:event-id/paged", organizerHandler.ListPaged)

		protectedOrganizers.POST("", organizerHandler.Create)
		protectedOrganizers.PUT("/:id", organizerHandler.Update)
		protectedOrganizers.DELETE("/:id", organizerHandler.Delete)
	}

	speakers := api.Group("/speakers")
	protectedSpeakers := speakers.Group("/")
	protectedSpeakers.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
	{
		speakers.GET("/event/:event-id", speakerHandler.GetAll)
		speakers.GET("/id/:id", speakerHandler.GetByID)
		speakers.GET("/event/:event-id/paged", speakerHandler.ListPaged)

		protectedSpeakers.POST("", speakerHandler.Create)
		protectedSpeakers.PUT("/:id", speakerHandler.Update)
		protectedSpeakers.DELETE("/:id", speakerHandler.Delete)
	}

	return r
}
