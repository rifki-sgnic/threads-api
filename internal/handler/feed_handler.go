package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
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
		if errors.Is(err, pagination.ErrInvalidCursor) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	items := make([]postWithStatsResponse, 0, len(page.Items))
	for i := range page.Items {
		items = append(items, toPostWithStatsResponse(&page.Items[i]))
	}

	var next *string
	if page.NextCursor != "" {
		next = &page.NextCursor
	}

	c.JSON(http.StatusOK, gin.H{"data": items, "next_cursor": next})
}
