package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

type FeedHandler struct {
	svc service.FeedService
}

func NewFeedHandler(svc service.FeedService) *FeedHandler {
	return &FeedHandler{svc: svc}
}

func (h *FeedHandler) Get(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))

	page, err := h.svc.GetFeed(c.GetUint(middleware.UserIDKey), c.Query("cursor"), limit)
	if err != nil {
		respondError(c, err)
		return
	}

	respondPostPage(c, page)
}
