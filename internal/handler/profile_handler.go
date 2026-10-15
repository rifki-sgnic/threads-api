package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/pagination"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

type ProfileHandler struct {
	svc service.ProfileService
}

func NewProfileHandler(svc service.ProfileService) *ProfileHandler {
	return &ProfileHandler{svc: svc}
}

type profileResponse struct {
	ID             uint      `json:"id"`
	Username       string    `json:"username"`
	Bio            string    `json:"bio"`
	FollowerCount  int64     `json:"follower_count"`
	FollowingCount int64     `json:"following_count"`
	PostCount      int64     `json:"post_count"`
	IsFollowing    bool      `json:"is_following"`
	CreatedAt      time.Time `json:"created_at"`
}

type meResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

type updateMeRequest struct {
	Bio *string `json:"bio" binding:"required,max=160"`
}

func (h *ProfileHandler) Get(c *gin.Context) {
	p, err := h.svc.GetProfile(c.GetUint(middleware.UserIDKey), c.Param("username"))
	if err != nil {
		respondProfileError(c, err)
		return
	}

	c.JSON(http.StatusOK, profileResponse{
		ID:             p.User.ID,
		Username:       p.User.Username,
		Bio:            p.User.Bio,
		FollowerCount:  p.FollowerCount,
		FollowingCount: p.FollowingCount,
		PostCount:      p.PostCount,
		IsFollowing:    p.IsFollowing,
		CreatedAt:      p.User.CreatedAt,
	})
}

func (h *ProfileHandler) Posts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.Query("limit"))

	page, err := h.svc.GetUserPosts(
		c.GetUint(middleware.UserIDKey),
		c.Param("username"),
		c.Query("cursor"),
		limit,
	)
	if err != nil {
		respondProfileError(c, err)
		return
	}

	respondPostPage(c, page)
}

func (h *ProfileHandler) UpdateMe(c *gin.Context) {
	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.UpdateBio(c.GetUint(middleware.UserIDKey), *req.Bio)
	if err != nil {
		respondProfileError(c, err)
		return
	}

	c.JSON(http.StatusOK, meResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
	})
}

func respondProfileError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, pagination.ErrInvalidCursor):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
