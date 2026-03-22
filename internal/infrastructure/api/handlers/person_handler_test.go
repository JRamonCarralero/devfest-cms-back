package handlers_test

import (
	"bytes"
	"devfest/internal/domain"
	"devfest/internal/domain/mocks"
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

func MockAuthMiddleware(userId uuid.UUID) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userID", userId)
		c.Next()
	}
}

func setupPersonRouter(mockUsecase *mocks.PersonUsecase, authUserId uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := handlers.NewPersonHandler(mockUsecase)

	auth := r.Group("/")
	auth.Use(MockAuthMiddleware(authUserId))
	{
		auth.POST("/persons", h.Create)
		auth.PUT("/persons/:id", h.Update)
		auth.DELETE("/persons/:id", h.Delete)
	}

	r.GET("/persons/:id", h.GetByID)
	r.GET("/persons/paged", h.GetPaged)

	return r
}

func TestPersonHandler(t *testing.T) {
	mockUsecase := new(mocks.PersonUsecase)
	authUserId := uuid.New()
	router := setupPersonRouter(mockUsecase, authUserId)

	t.Run("Create - Should return 201", func(t *testing.T) {
		body := map[string]interface{}{
			"first_name": "Jane",
			"last_name":  "Doe",
			"email":      "jane@doe.com",
		}
		jsonBody, _ := json.Marshal(body)

		mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("*domain.Person")).
			Return(&domain.Person{ID: uuid.New(), FirstName: "Jane"}, nil).Once()

		req, _ := http.NewRequest(http.MethodPost, "/persons", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("GetByID - Should return 200", func(t *testing.T) {
		id := uuid.New()
		mockUsecase.On("GetByID", mock.Anything, id).
			Return(&domain.Person{ID: id, FirstName: "Jane"}, nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/persons/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("ListPaged - Should return 200 with data", func(t *testing.T) {
		mockUsecase.On("ListPaged", mock.Anything, "search", int32(1), int32(10)).
			Return([]domain.Person{}, int64(0), nil).Once()

		req, _ := http.NewRequest(http.MethodGet, "/persons/paged?search=search&page=1&pageSize=10", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Update - Should return 200", func(t *testing.T) {
		id := uuid.New()
		body := map[string]interface{}{"first_name": "Updated"}
		jsonBody, _ := json.Marshal(body)

		mockUsecase.On("Update", mock.Anything, id, mock.AnythingOfType("*domain.UpdatePerson")).
			Return(&domain.Person{ID: id, FirstName: "Updated"}, nil).Once()

		req, _ := http.NewRequest(http.MethodPut, "/persons/"+id.String(), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Delete - Should return 204", func(t *testing.T) {
		id := uuid.New()
		mockUsecase.On("Delete", mock.Anything, id).Return(nil).Once()

		req, _ := http.NewRequest(http.MethodDelete, "/persons/"+id.String(), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
	})

	t.Run("Error Case - Invalid UUID returns 400", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/persons/not-a-uuid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
