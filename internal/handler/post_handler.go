package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/models"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

type PostHandler struct {
	svc service.PostService
}

func NewPostHandler(svc service.PostService) *PostHandler {
	return &PostHandler{svc: svc}
}

type createPostRequest struct {
	Content  string `json:"content" binding:"required,max=500"`
	ParentID *uint  `json:"parent_id"`
}

type authorResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type postResponse struct {
	ID        uint           `json:"id"`
	Content   string         `json:"content"`
	ParentID  *uint          `json:"parent_id"`
	Author    authorResponse `json:"author"`
	CreatedAt time.Time      `json:"created_at"`
}

type postWithStatsResponse struct {
	postResponse
	LikeCount  int64 `json:"like_count"`
	ReplyCount int64 `json:"reply_count"`
	LikedByMe  bool  `json:"liked_by_me"`
}

func toPostResponse(p *models.Post) postResponse {
	return postResponse{
		ID:        p.ID,
		Content:   p.Content,
		ParentID:  p.ParentID,
		Author:    authorResponse{ID: p.UserID, Username: p.User.Username},
		CreatedAt: p.CreatedAt,
	}
}

func toPostWithStatsResponse(v *service.PostView) postWithStatsResponse {
	return postWithStatsResponse{
		postResponse: toPostResponse(&v.Post),
		LikeCount:    v.LikeCount,
		ReplyCount:   v.ReplyCount,
		LikedByMe:    v.LikedByMe,
	}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req createPostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post, err := h.svc.Create(service.CreatePostInput{
		UserID:   c.GetUint(middleware.UserIDKey),
		Content:  req.Content,
		ParentID: req.ParentID,
	})

	if err != nil {
		respondPostError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toPostResponse(post))
}

func (h *PostHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	view, err := h.svc.GetDetail(c.GetUint(middleware.UserIDKey), id)
	if err != nil {
		respondPostError(c, err)
		return
	}

	c.JSON(http.StatusOK, toPostWithStatsResponse(view))
}

func (h *PostHandler) Replies(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	views, err := h.svc.GetRepliesWithStats(c.GetUint(middleware.UserIDKey), id)
	if err != nil {
		respondPostError(c, err)
		return
	}

	out := make([]postWithStatsResponse, 0, len(views))
	for i := range views {
		out = append(out, toPostWithStatsResponse(&views[i]))
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

func (h *PostHandler) Delete(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}

	if err := h.svc.Delete(c.GetUint(middleware.UserIDKey), id); err != nil {
		respondPostError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func respondPostError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrPostNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrEmptyContent):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
