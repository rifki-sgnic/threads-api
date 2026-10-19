package handler

import (
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
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *FollowHandler) Unfollow(c *gin.Context) {
	err := h.svc.Unfollow(c.GetUint(middleware.UserIDKey), c.Param("username"))
	if err != nil {
		respondError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
