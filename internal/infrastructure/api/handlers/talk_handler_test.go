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

func setupTalkRouter(mockUC *mocks.TalkUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewTalkHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/talks", h.Create)
		auth.PUT("/talks/:id", h.Update)
		auth.DELETE("/talks/:id", h.Delete)
		auth.POST("/talks/speakers", h.AddSpeaker)
		auth.DELETE("/talks/speakers", h.RemoveSpeaker)
	}

	r.GET("/events/:event-id/talks", h.GetAll)
	r.GET("/talks/:id", h.GetByID)

	return r
}

func TestTalkHandler(t *testing.T) {
	mockUC := new(mocks.TalkUsecase)
	authUserID := uuid.New()
	router := setupTalkRouter(mockUC, authUserID)
	ctxMatch := mock.Anything

	t.Run("Create - Success", func(t *testing.T) {
		eventID := uuid.New()
		tags := []string{"Go", "Backend"}
		body := map[string]interface{}{
			"event_id":    eventID.String(),
			"title":       "Testing Talk",
			"description": "A great talk",
			"tags":        tags,
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", ctxMatch, mock.MatchedBy(func(tk *domain.Talk) bool {
			return tk.Title == "Testing Talk" && tk.CreatedBy == authUserID
		})).Return(&domain.Talk{ID: uuid.New(), Title: "Testing Talk"}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/talks", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetAll - Success", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", ctxMatch, eventID).
			Return([]domain.Talk{{ID: uuid.New(), Title: "Talk 1"}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/events/%s/talks", eventID), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Success", func(t *testing.T) {
		id := uuid.New()
		title := "Updated Title"
		body := map[string]interface{}{"title": title}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", ctxMatch, id, mock.MatchedBy(func(u *domain.UpdateTalk) bool {
			return *u.Title == title && u.UpdatedBy == authUserID
		})).Return(&domain.Talk{ID: id, Title: title}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/talks/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("AddSpeaker - Success", func(t *testing.T) {
		talkID := uuid.New()
		speakerID := uuid.New()
		body := map[string]interface{}{
			"talk_id":    talkID.String(),
			"speaker_id": speakerID.String(),
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("AddSpeaker", ctxMatch, mock.MatchedBy(func(ts *domain.TalkSpeaker) bool {
			return ts.TalkID == talkID && ts.SpeakerID == speakerID && ts.CreatedBy == authUserID
		})).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/talks/speakers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("RemoveSpeaker - Success (Query Params)", func(t *testing.T) {
		talkID := uuid.New()
		speakerID := uuid.New()

		mockUC.On("RemoveSpeaker", ctxMatch, mock.MatchedBy(func(ts *domain.TalkSpeaker) bool {
			return ts.TalkID == talkID && ts.SpeakerID == speakerID
		})).Return(nil).Once()

		url := fmt.Sprintf("/talks/speakers?talk_id=%s&speaker_id=%s", talkID, speakerID)
		req, _ := http.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("GetByID - Fail Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/talks/not-a-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
