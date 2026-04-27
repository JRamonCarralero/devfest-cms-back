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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTrackRouter(mockUC *mocks.TrackUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewTrackHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/tracks", h.Create)
		auth.PUT("/tracks/:id", h.Update)
		auth.DELETE("/tracks/:id", h.Delete)
	}

	r.GET("/events/:event-id/tracks", h.GetAll)
	r.GET("/events/:event-id/schedule", h.GetFullEventSchedule)
	r.GET("/tracks/:id", h.GetByID)

	return r
}

func TestTrackHandler(t *testing.T) {
	mockUC := new(mocks.TrackUsecase)
	authUserID := uuid.New()
	router := setupTrackRouter(mockUC, authUserID)
	ctx := mock.Anything

	t.Run("Create - Should return 201", func(t *testing.T) {
		eventID := uuid.New()
		now := time.Now()
		body := map[string]interface{}{
			"event_id":   eventID.String(),
			"name":       "Main Track",
			"event_date": now.Format(time.RFC3339Nano),
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", ctx, mock.MatchedBy(func(tr *domain.Track) bool {
			return tr.Name == "Main Track" && tr.EventID == eventID && tr.CreatedBy == authUserID
		})).Return(&domain.Track{ID: uuid.New(), Name: "Main Track", EventID: eventID}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/tracks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetAll - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", ctx, eventID).
			Return([]domain.Track{{ID: uuid.New(), Name: "Track 1"}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/tracks", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dtos.TrackResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Len(t, resp, 1)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetById", ctx, id).
			Return(&domain.Track{ID: id, Name: "Track 1"}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/tracks/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		newName := "Updated Name"
		body := map[string]interface{}{"name": newName}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", ctx, id, mock.MatchedBy(func(u *domain.UpdateTrack) bool {
			return u.Name != nil && *u.Name == newName && u.UpdatedBy == authUserID
		})).Return(&domain.Track{ID: id, Name: newName}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/tracks/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", ctx, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/tracks/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("GetFullEventSchedule - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetFullEventSchedule", ctx, eventID).
			Return([]domain.FullTrackSchedule{
				{TrackID: uuid.New(), TrackName: "Main"},
			}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/schedule", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp []dtos.FullTrackScheduleResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NotEmpty(t, resp)
	})

	t.Run("Invalid UUID - Should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/tracks/invalid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
