package response

import (
	"devfest/internal/domain"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	traceID, exists := c.Get("trace_id")
	sTraceID := "unknown"
	if exists {
		sTraceID = traceID.(string)
	}

	var appErr *domain.AppError
	if !errors.As(err, &appErr) {
		appErr = domain.NewAppError(domain.TypeInternal, "An unexpected error occurred", err)
	}

	appErr.TraceID = sTraceID

	statusCode := http.StatusInternalServerError
	switch appErr.Type {
	case domain.TypeNotFound:
		statusCode = http.StatusNotFound
	case domain.TypeBadRequest:
		statusCode = http.StatusBadRequest
	case domain.TypeAlreadyExists:
		statusCode = http.StatusConflict
	case domain.TypeUnauthenticated:
		statusCode = http.StatusUnauthorized
	}

	log.Printf("[ERROR] [%s] TraceID: %s | Location: %s | Msg: %s | InternalErr: %v",
		appErr.Type, appErr.TraceID, appErr.Location, appErr.Message, appErr.Err)

	c.JSON(statusCode, gin.H{
		"error":    appErr.Message,
		"type":     appErr.Type,
		"trace_id": appErr.TraceID,
	})
}
