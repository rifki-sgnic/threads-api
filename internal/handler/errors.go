package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPostNotFound),
		errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})

	case errors.Is(err, service.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})

	case errors.Is(err, service.ErrUserExists):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, service.ErrEmptyContent),
		errors.Is(err, service.ErrCannotFollowSelf),
		errors.Is(err, pagination.ErrInvalidCursor):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	default:
		middleware.Logger(c).Error("unhandled error", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
