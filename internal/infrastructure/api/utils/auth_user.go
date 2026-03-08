package utils

import (
	"devfest/internal/domain"
	"devfest/internal/infrastructure/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserId(c *gin.Context) (uuid.UUID, error) {
	userID, exists := c.Get(middleware.CtxUserID)
	if !exists {
		return uuid.Nil, domain.NewAppError(domain.TypeUnauthorized, "User context missing", nil)
	}

	uid, ok := userID.(uuid.UUID)
	if !ok {
		return uuid.Nil, domain.NewAppError(domain.TypeInternal, "User ID in context is not a UUID", nil)
	}

	return uid, nil
}

func GetUserRole(c *gin.Context) (domain.Role, error) {
	userRole, exists := c.Get(middleware.CtxRole)
	if !exists {
		return domain.Role(""), domain.NewAppError(domain.TypeUnauthorized, "User context missing", nil)
	}

	role, ok := userRole.(domain.Role)
	if !ok {
		return domain.Role(""), domain.NewAppError(domain.TypeInternal, "User role in context is not a Role", nil)
	}

	return role, nil
}
