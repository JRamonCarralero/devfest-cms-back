package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseIntOrDefault parses a string to an int and returns the default value if the parsing fails
func ParseIntOrDefault(value string, defaultValue int, min int) int {
	if value == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(value)
	if err != nil || i < min {
		return defaultValue
	}

	return i
}

// GetPaginationParams returns the pagination parameters from the request
func GetPaginationParams(c *gin.Context) (page, pageSize int) {
	page = ParseIntOrDefault(c.Query("page"), 1, 1)
	pageSize = ParseIntOrDefault(c.Query("pageSize"), 10, 1)
	return
}
