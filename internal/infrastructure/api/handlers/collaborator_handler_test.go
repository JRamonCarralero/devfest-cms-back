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

func setupCollaboratorRouter(mockUC *mocks.CollaboratorUsecase, authUserID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewCollaboratorHandler(mockUC)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserID))
	{
		auth.POST("/collaborators", h.Create)
		auth.PUT("/collaborators/:id", h.Update)
		auth.DELETE("/collaborators/:id", h.Delete)
	}

	r.GET("/events/:event-id/collaborators", h.GetAll)
	r.GET("/events/:event-id/collaborators/paged", h.ListPaged)
	r.GET("/collaborators/:id", h.GetByID)

	return r
}

func TestCollaboratorHandler(t *testing.T) {
	mockUC := new(mocks.CollaboratorUsecase)
	authUserID := uuid.New()
	router := setupCollaboratorRouter(mockUC, authUserID)

	t.Run("Create - Should return 201", func(t *testing.T) {
		eventID := uuid.New()
		personID := uuid.New()
		body := map[string]interface{}{
			"event_id":  eventID.String(),
			"person_id": personID.String(),
			"area":      "Marketing",
		}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Create", mock.Anything, mock.MatchedBy(func(c *domain.Collaborator) bool {
			return c.EventID == eventID && c.Person.ID == personID && c.CreatedBy == authUserID
		})).Return(&domain.Collaborator{ID: uuid.New(), EventID: eventID}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/collaborators", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockUC.AssertExpectations(t)
	})

	t.Run("GetAll - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("GetAll", mock.Anything, eventID).
			Return([]domain.Collaborator{{ID: uuid.New(), EventID: eventID}}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/collaborators", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		area := "Logistics"
		body := map[string]interface{}{"area": area}
		jsonBody, _ := json.Marshal(body)

		mockUC.On("Update", mock.Anything, id, mock.MatchedBy(func(u *domain.UpdateCollaborator) bool {
			return *u.Area == area && u.UpdatedBy == authUserID
		})).Return(&domain.Collaborator{ID: id, Area: &area}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/collaborators/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("Delete", mock.Anything, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/collaborators/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("ListPaged - Should return 200", func(t *testing.T) {
		eventID := uuid.New()
		mockUC.On("ListPaged", mock.Anything, eventID, "test", int32(1), int32(10)).
			Return([]domain.Collaborator{}, int64(0), nil).Once()

		url := "/events/" + eventID.String() + "/collaborators/paged?search=test&page=1&pageSize=10"
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUC.On("GetById", mock.Anything, id).
			Return(&domain.Collaborator{
				ID:     id,
				Person: domain.Person{ID: uuid.New(), FirstName: "John"},
			}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/collaborators/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp dtos.CollaboratorResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, id, resp.ID)
	})

	t.Run("GetByID - Invalid UUID should return 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/collaborators/esto-no-es-un-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
