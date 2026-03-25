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

func setupDeveloperRouter(mockUC *mocks.DeveloperUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewDeveloperHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/developers", h.Create)
		auth.PUT("/developers/:id", h.Update)
		auth.DELETE("/developers/:id", h.Delete)
	}

	r.GET("/events/:event-id/developers", h.GetAll)
	r.GET("/events/:event-id/developers/paged", h.ListPaged)
	r.GET("/developers/:id", h.GetByID)

	return r
}

func TestDeveloperHandler(t *testing.T) {
	mockUC := new(mocks.DeveloperUsecase)
	authUserID := uuid.New()
	router := setupDeveloperRouter(mockUC, authUserID)

	t.Run("Create - Should return 201", func(t *testing.T) {
		eventID := uuid.New()
		personID := uuid.New()
		body := map[string]interface{}{
			"event_id":  eventID.String(),
			"person_id": personID.String(),
			"area":      "Marketing",
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Developer) bool {
			return c.EventID == eventID && c.Person.ID == personID && c.CreatedBy == authUserID
		})).Return(&domain.Developer{ID: uuid.New(), EventID: eventID}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/developers", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetAll - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", mock.Anything, eventID).
			Return([]domain.Developer{{ID: uuid.New(), EventID: eventID}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/developers", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		roleDescription := "Backend"
		body := map[string]interface{}{"role_description": roleDescription}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", mock.Anything, id, mock.MatchedBy(func(u *domain.UpdateDeveloper) bool {
			return *u.RoleDescription == roleDescription && u.UpdatedBy == authUserID
		})).Return(&domain.Developer{ID: id, RoleDescription: &roleDescription}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/developers/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", mock.Anything, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/developers/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ListPaged - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("ListPaged", mock.Anything, eventID, "test", int32(1), int32(10)).
			Return([]domain.Developer{}, int64(0), nil).Once()

		url := "/events/" + eventID.String() + "/developers/paged?search=test&page=1&pageSize=10"
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetById", mock.Anything, id).
			Return(&domain.Developer{
				ID:     id,
				Person: domain.Person{ID: uuid.New(), FirstName: "John"},
			}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/developers/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp dtos.DeveloperResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("GetByID - Invalid UUID should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/developers/esto-no-es-un-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
