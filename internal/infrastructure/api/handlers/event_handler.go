package handlers

import (
	"devfest/internal/domain"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	usecase domain.EventUsecase
}

func NewEventHandler(usecase domain.EventUsecase) *EventHandler {
	return &EventHandler{usecase: usecase}
}

// GetEvents
func (h *EventHandler) GetEvents(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	orderBy := c.DefaultQuery("order", "created_at_desc")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	events, total, err := h.usecase.GetEvents(c.Request.Context(), search, int32(page), int32(pageSize), orderBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch events"})
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

// GetBySlug
func (h *EventHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	event, err := h.usecase.GetEventBySlug(c.Request.Context(), slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}
