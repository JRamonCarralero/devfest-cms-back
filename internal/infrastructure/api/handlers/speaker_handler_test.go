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

func setupSpeakerRouter(mockUC *mocks.SpeakerUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewSpeakerHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/speakers", h.Create)
		auth.PUT("/speakers/:id", h.Update)
		auth.DELETE("/speakers/:id", h.Delete)
	}

	r.GET("/events/:event-id/speakers", h.GetAll)
	r.GET("/events/:event-id/speakers/paged", h.ListPaged)
	r.GET("/speakers/:id", h.GetByID)

	return r
}

func TestSpeakerHandler(t *testing.T) {
	mockUC := new(mocks.SpeakerUsecase)
	authUserID := uuid.New()
	router := setupSpeakerRouter(mockUC, authUserID)

	t.Run("Create - Should return 201", func(t *testing.T) {
		eventID := uuid.New()
		personID := uuid.New()
		body := map[string]interface{}{
			"event_id":  eventID.String(),
			"person_id": personID.String(),
			"area":      "Marketing",
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Speaker) bool {
			return c.EventID == eventID && c.Person.ID == personID && c.CreatedBy == authUserID
		})).Return(&domain.Speaker{ID: uuid.New(), EventID: eventID}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/speakers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetAll - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", mock.Anything, eventID).
			Return([]domain.Speaker{{ID: uuid.New(), EventID: eventID}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/speakers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		bio := "AI Expert"
		body := map[string]interface{}{"bio": bio}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", mock.Anything, id, mock.MatchedBy(func(u *domain.UpdateSpeaker) bool {
			return *u.Bio == bio && u.UpdatedBy == authUserID
		})).Return(&domain.Speaker{ID: id, Bio: &bio}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/speakers/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", mock.Anything, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/speakers/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ListPaged - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("ListPaged", mock.Anything, eventID, "test", int32(1), int32(10)).
			Return([]domain.Speaker{}, int64(0), nil).Once()

		url := "/events/" + eventID.String() + "/speakers/paged?search=test&page=1&pageSize=10"
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetById", mock.Anything, id).
			Return(&domain.Speaker{
				ID:     id,
				Person: domain.Person{ID: uuid.New(), FirstName: "John"},
			}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/speakers/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp dtos.SpeakerResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("GetByID - Invalid UUID should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/speakers/esto-no-es-un-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
