package middleware_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware_Roles(t *testing.T) {
	gin.SetMode(gin.TestMode)

	userID := uuid.New().String()

	t.Run("Access Granted - Admin accessing Admin route", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.AuthMiddleware(domain.RoleAdmin, domain.RoleSuperAdmin))
		r.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		token := generateTestToken(userID, string(domain.RoleAdmin))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("Access Denied - User accessing Admin route", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.AuthMiddleware(domain.RoleAdmin))
		r.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		token := generateTestToken(userID, "user")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})

	t.Run("Invalid UUID in token", func(t *testing.T) {
		r := gin.New()
		r.Use(middleware.AuthMiddleware(domain.RoleAdmin))
		r.GET("/protected", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		token := generateTestToken("this-is-not-a-uuid", "admin")

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/protected", nil)
		req.Header.Set("Authorization", "Bearer "+token)

		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}

func TestTraceMiddleware(t *testing.T) {
	r := gin.New()
	r.Use(middleware.TraceMiddleware())
	r.GET("/trace", func(c *gin.Context) {
		val, exists := c.Get("trace_id")
		assert.True(t, exists)
		assert.NotEmpty(t, val)
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/trace", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Trace-ID"))
}

func generateTestToken(userID string, role string) string {
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))

	payloadMap := map[string]interface{}{
		"sub": userID,
		"app_metadata": map[string]interface{}{
			"role": role,
		},
	}
	payloadBytes, _ := json.Marshal(payloadMap)
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)

	signature := base64.RawURLEncoding.EncodeToString([]byte("signature-data"))

	return fmt.Sprintf("%s.%s.%s", header, payload, signature)
}
