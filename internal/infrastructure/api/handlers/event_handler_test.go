package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/infrastructure/api/handlers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEventHandler_All(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx := mock.Anything

	t.Run("GetEvents", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		events := []domain.Event{{Name: "E1"}, {Name: "E2"}}

		mockUC.On("GetEvents", mock.Anything).Return(events, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/events", nil)

		handler.GetEvents(c)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "E1")
	})

	t.Run("GetByID", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		id := uuid.New()

		mockUC.On("GetByID", ctx, id).Return(&domain.Event{ID: id, Name: "Test"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: id.String()}}
		c.Request, _ = http.NewRequest("GET", "/", nil)

		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetBySlug", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		slug := "devfest-2026"

		mockUC.On("GetEventBySlug", ctx, slug).Return(&domain.Event{Slug: slug}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "slug", Value: slug}}
		c.Request, _ = http.NewRequest("GET", "/", nil)

		handler.GetBySlug(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetActive", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)

		mockUC.On("GetActiveEvents", ctx).Return([]domain.Event{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("GET", "/active", nil)

		handler.GetActive(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetPaged", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)

		mockUC.On("GetEventsPaged", mock.Anything, "search", int32(1), int32(10), "created_at_desc").
			Return([]domain.Event{}, int64(0), nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		req, _ := http.NewRequest("GET", "/events?search=search&page=1&pageSize=10", nil)
		c.Request = req

		handler.GetPaged(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Create", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		userID := uuid.New()

		body := map[string]interface{}{"name": "New Event", "slug": "new-event", "is_active": true}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("CreateEvent", mock.Anything, mock.AnythingOfType("*domain.Event")).
			Return(&domain.Event{Name: "New Event"}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest("POST", "/", bytes.NewBuffer(jsonBody))
		c.Request.Header.Set("Content-Type", "application/json")

		c.Set("userID", userID)

		handler.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Update", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		eventID := uuid.New()
		userID := uuid.New()

		name := "Updated Name"
		body := map[string]interface{}{"name": name}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("UpdateEvent", mock.Anything, eventID, mock.AnythingOfType("*domain.UpdateEvent")).
			Return(&domain.Event{ID: eventID, Name: name}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "id", Value: eventID.String()}}
		c.Request, _ = http.NewRequest("PUT", "/", bytes.NewBuffer(jsonBody))

		c.Set("userID", userID)

		handler.Update(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete", func(t *testing.T) {
		mockUC := mocks.NewEventUsecase(t)
		handler := handlers.NewEventHandler(mockUC)
		eventID := uuid.New()

		mockUC.On("DeleteEvent", mock.Anything, eventID).Return(nil).Once()

		r := gin.New()
		r.DELETE("/events/:id", handler.Delete)

		req, _ := http.NewRequest("DELETE", "/events/"+eventID.String(), nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockUC.AssertExpectations(t)
	})
}
