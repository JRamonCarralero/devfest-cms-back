package handlers_test

import (
	"bytes"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
	"devfest/internal/infrastructure/api/handlers"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSchedulerRouter(mockUC *mocks.SchedulerUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewSchedulerHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/scheduler", h.Create)
		auth.PUT("/scheduler/:id", h.Update)
		auth.DELETE("/scheduler/:id", h.Delete)
	}

	r.GET("/tracks/:track-id/scheduler", h.GetAllByTrack)
	r.GET("/scheduler/:id", h.GetByID)

	return r
}

func TestSchedulerHandler(t *testing.T) {
	mockUC := new(mocks.SchedulerUsecase)
	authUserID := uuid.New()
	router := setupSchedulerRouter(mockUC, authUserID)
	ctxMatch := mock.Anything

	t.Run("Create - Success", func(t *testing.T) {
		trackID := uuid.New()
		talkID := uuid.New()
		room := "Main Hall"
		startTimeStr := "2026-10-10T10:00:00Z"
		endTimeStr := "2026-10-10T11:00:00Z"

		body := map[string]interface{}{
			"track_id":   trackID.String(),
			"talk_id":    talkID.String(),
			"start_time": startTimeStr,
			"end_time":   endTimeStr,
			"room":       room,
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", ctxMatch, mock.MatchedBy(func(s *domain.Scheduler) bool {
			return s.Track.ID == trackID && s.Room == room
		})).Return(&domain.Scheduler{ID: uuid.New(), Room: room}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/scheduler", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("Create - Fail Invalid Date Format", func(t *testing.T) {
		body := map[string]interface{}{
			"track_id":   uuid.New().String(),
			"talk_id":    uuid.New().String(),
			"start_time": "10-10-2026 10:00",
			"end_time":   "2026-10-10T11:00:00Z",
			"room":       "Room A",
		}
		jsonBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPost, "/scheduler", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetAllByTrack - Success", func(t *testing.T) {
		trackID := uuid.New()
		mockUC.On("GetAllByTrack", ctxMatch, trackID).
			Return([]domain.Scheduler{{ID: uuid.New(), Room: "Room B"}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/tracks/%s/scheduler", trackID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID - Success", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetByID", ctxMatch, id).
			Return(&domain.Scheduler{ID: id, Room: "Room C"}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/scheduler/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Success", func(t *testing.T) {
		id := uuid.New()
		newRoom := "Workshop Room"
		body := map[string]interface{}{"room": newRoom}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", ctxMatch, id, mock.MatchedBy(func(u *domain.UpdateScheduler) bool {
			return *u.Room == newRoom && u.UpdatedBy == authUserID
		})).Return(&domain.Scheduler{ID: id, Room: newRoom}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/scheduler/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Success", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", ctxMatch, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/scheduler/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Invalid UUID Parameter - Should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/scheduler/not-a-valid-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
