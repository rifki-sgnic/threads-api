package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

type LikeHandler struct {
	svc service.LikeService
}

func NewLikeHandler(svc service.LikeService) *LikeHandler {
	return &LikeHandler{svc: svc}
}

func (h *LikeHandler) Like(c *gin.Context) {
	postID, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.Like(c.GetUint(middleware.UserIDKey), postID); err != nil {
		respondPostError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *LikeHandler) Unlike(c *gin.Context) {
	postID, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.Unlike(c.GetUint(middleware.UserIDKey), postID); err != nil {
		respondPostError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}
