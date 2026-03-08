package handlers

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/infrastructure/api/response"
	"net/http"
	"strconv"

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

	c.JSON(http.StatusOK, events)
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

	c.JSON(http.StatusOK, event)
}

// GetBySlug
func (h *EventHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	event, err := h.usecase.GetEventBySlug(c.Request.Context(), slug)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetActive
func (h *EventHandler) GetActive(c *gin.Context) {
	events, err := h.usecase.GetActiveEvents(c.Request.Context())
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetPaged
func (h *EventHandler) GetPaged(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	orderBy := c.DefaultQuery("order", "created_at_desc")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	events, total, err := h.usecase.GetEventsPaged(c.Request.Context(), search, int32(page), int32(pageSize), orderBy)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": events,
		"meta": gin.H{
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
		},
	})
}

// Create
func (h *EventHandler) Create(c *gin.Context) {
	var dto dtos.CreateEventDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		newErr := domain.NewAppError(domain.TypeBadRequest, "Invalid request body", err)
		response.HandleError(c, newErr)
		return
	}

	event, err := h.usecase.CreateEvent(c.Request.Context(), dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
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

	event, err := h.usecase.UpdateEvent(c.Request.Context(), id, dto)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
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
