package handlers

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/infrastructure/api/response"
	"devfest/internal/infrastructure/api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	usecase domain.EventUsecase
}

func NewEventHandler(usecase domain.EventUsecase) *EventHandler {
	return &EventHandler{usecase: usecase}
}

// GetEvents
func (h *EventHandler) GetEvents(c *gin.Context) {
	events, err := h.usecase.GetEvents(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	if len(events) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	eventsResponse := make([]dtos.EventResponse, len(events))
	for i, event := range events {
		eventsResponse[i] = mapToDomainEventResponse(event)
	}

	c.JSON(http.StatusOK, eventsResponse)
}

// GetByID
func (h *EventHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(c, newErr)
		return
	}

	event, err := h.usecase.GetByID(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	eventResponse := mapToDomainEventResponse(*event)

	c.JSON(http.StatusOK, eventResponse)
}

// GetBySlug
func (h *EventHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	event, err := h.usecase.GetEventBySlug(c.Request.Context(), slug)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	eventResponse := mapToDomainEventResponse(*event)

	c.JSON(http.StatusOK, eventResponse)
}

// GetActive
func (h *EventHandler) GetActive(c *gin.Context) {
	events, err := h.usecase.GetActiveEvents(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	eventsResponse := make([]dtos.EventResponse, len(events))
	for i, event := range events {
		eventsResponse[i] = mapToDomainEventResponse(event)
	}

	c.JSON(http.StatusOK, eventsResponse)
}

// GetPaged
func (h *EventHandler) GetPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	orderBy := c.DefaultQuery("order", "created_at_desc")

	page, pageSize := utils.GetPaginationParams(c)
	if pageSize > 100 {
		pageSize = 100
	}

	events, total, err := h.usecase.GetEventsPaged(c.Request.Context(), search, int32(page), int32(pageSize), orderBy)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	if len(events) == 0 {
		c.Status(http.StatusNoContent)
		return
	}

	eventsResponse := make([]dtos.EventResponse, len(events))
	for i, event := range events {
		eventsResponse[i] = mapToDomainEventResponse(event)
	}

	res := dtos.PagedResponse[dtos.EventResponse]{
		Data: eventsResponse,
		Meta: dtos.PagedMeta{
			Total:    total,
			Page:     int32(page),
			PageSize: int32(pageSize),
		},
	}

	c.JSON(http.StatusOK, res)
}

// Create
func (h *EventHandler) Create(c *gin.Context) {
	var dto dtos.CreateEventDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid request body", err)
		response.HandleError(c, newErr)
		return
	}

	uid, err := utils.GetUserId(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	dto.CreatedBy = uid

	isActive := false
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}
	dto.IsActive = &isActive

	newEvent := domain.Event{
		Name:     dto.Name,
		Slug:     dto.Slug,
		IsActive: dto.IsActive,
		Audit: domain.Audit{
			CreatedBy: dto.CreatedBy,
			UpdatedBy: dto.CreatedBy,
		},
	}

	event, err := h.usecase.CreateEvent(c.Request.Context(), &newEvent)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	eventsResponse := mapToDomainEventResponse(*event)

	c.JSON(http.StatusCreated, eventsResponse)
}

// Update
func (h *EventHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(c, newErr)
		return
	}

	var dto dtos.UpdateEventDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid request body", err)
		response.HandleError(c, newErr)
		return
	}

	uid, err := utils.GetUserId(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	upEvent := domain.UpdateEvent{
		Name:      dto.Name,
		Slug:      dto.Slug,
		IsActive:  dto.IsActive,
		UpdatedBy: uid,
	}

	event, err := h.usecase.UpdateEvent(c.Request.Context(), id, &upEvent)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	eventsResponse := mapToDomainEventResponse(*event)

	c.JSON(http.StatusOK, eventsResponse)
}

// Delete
func (h *EventHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid ID", err)
		response.HandleError(c, newErr)
		return
	}

	err = h.usecase.DeleteEvent(c.Request.Context(), id)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func mapToDomainEventResponse(event domain.Event) dtos.EventResponse {
	return dtos.EventResponse{
		ID:        event.ID,
		Name:      event.Name,
		Slug:      event.Slug,
		IsActive:  event.IsActive,
		CreatedAt: event.Audit.CreatedAt,
		UpdatedAt: event.Audit.UpdatedAt,
		CreatedBy: event.Audit.CreatedBy,
		UpdatedBy: event.Audit.UpdatedBy,
	}
}
