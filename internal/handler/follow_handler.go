package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

type FollowHandler struct {
	svc service.FollowService
}

func NewFollowHandler(svc service.FollowService) *FollowHandler {
	return &FollowHandler{svc: svc}
}

func (h *FollowHandler) Follow(c *gin.Context) {
	err := h.svc.Follow(c.GetUint(middleware.UserIDKey), c.Param("username"))
	if err != nil {
		respondFollowError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	err := h.svc.Unfollow(c.GetUint(middleware.UserIDKey), c.Param("username"))
	if err != nil {
		respondFollowError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func respondFollowError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrCannotFollowSelf):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
