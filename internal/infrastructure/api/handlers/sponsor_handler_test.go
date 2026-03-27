package handlers_test

import (
	"bytes"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/infrastructure/api/dtos"
	"devfest/internal/infrastructure/api/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSponsorRouter(mockUC *mocks.SponsorUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewSponsorHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/Sponsors", h.Create)
		auth.PUT("/Sponsors/:id", h.Update)
		auth.DELETE("/Sponsors/:id", h.Delete)
	}

	r.GET("/events/:event-id/Sponsors", h.GetAll)
	r.GET("/events/:event-id/Sponsors/paged", h.ListPaged)
	r.GET("/Sponsors/:id", h.GetByID)

	return r
}

func TestSponsorHandler(t *testing.T) {
	mockUC := new(mocks.SponsorUsecase)
	authUserID := uuid.New()
	router := setupSponsorRouter(mockUC, authUserID)

	t.Run("Create - Should return 201", func(t *testing.T) {
		eventID := uuid.New()
		body := map[string]interface{}{
			"event_id":    eventID.String(),
			"name":        "Google",
			"logo_url":    "http:/www.google.com",
			"website_url": "http:/www.google.com",
			"tier":        "Gold",
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Sponsor) bool {
			return c.EventID == eventID && c.CreatedBy == authUserID
		})).Return(&domain.Sponsor{ID: uuid.New(), EventID: eventID}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/Sponsors", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetAll - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", mock.Anything, eventID).
			Return([]domain.Sponsor{{ID: uuid.New(), EventID: eventID}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/Sponsors", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		name := "Google Devs"
		body := map[string]interface{}{"name": name}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", mock.Anything, id, mock.MatchedBy(func(u *domain.UpdateSponsor) bool {
			return *u.Name == name && u.UpdatedBy == authUserID
		})).Return(&domain.Sponsor{ID: id, Name: name}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/Sponsors/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", mock.Anything, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/Sponsors/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ListPaged - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("ListPaged", mock.Anything, eventID, "test", int32(1), int32(10)).
			Return([]domain.Sponsor{}, int64(0), nil).Once()

		url := "/events/" + eventID.String() + "/Sponsors/paged?search=test&page=1&pageSize=10"
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetById", mock.Anything, id).
			Return(&domain.Sponsor{
				ID:   id,
				Name: "Google",
			}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/Sponsors/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp dtos.SponsorResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("GetByID - Invalid UUID should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/Sponsors/esto-no-es-un-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
