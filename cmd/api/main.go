package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/config"
	"github.com/rifki-sgnic/threads-api/internal/database"
	"github.com/rifki-sgnic/threads-api/internal/handler"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}

	// repos
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	likeRepo := repository.NewLikeRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// services
	statsLoader := service.NewPostStatsLoader(postRepo, likeRepo)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, time.Duration(cfg.JWTExpiryHours)*time.Hour)
	postSvc := service.NewPostService(postRepo, statsLoader)
	likeSvc := service.NewLikeService(likeRepo, postRepo)
	followSvc := service.NewFollowService(followRepo, userRepo)
	feedSvc := service.NewFeedService(postRepo, statsLoader)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)
	postHandler := handler.NewPostHandler(postSvc)
	likeHandler := handler.NewLikeHandler(likeSvc)
	followHandler := handler.NewFollowHandler(followSvc)
	feedHandler := handler.NewFeedHandler(feedSvc)

	authMW := middleware.Auth(cfg.JWTSecret)
	optAuthMW := middleware.OptionalAuth(cfg.JWTSecret)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group("/auth")
	auth.POST("/register", authHandler.Register)
	auth.POST("/login", authHandler.Login)

	// posts
	posts := r.Group("/posts")
	posts.GET("/:id", optAuthMW, postHandler.Get)
	posts.GET("/:id/replies", optAuthMW, postHandler.Replies)
	posts.POST("", authMW, postHandler.Create)
	posts.DELETE("/:id", authMW, postHandler.Delete)
	// like post
	posts.POST("/:id/like", authMW, likeHandler.Like)
	posts.DELETE("/:id/like", authMW, likeHandler.Unlike)

	// follow
	users := r.Group("/users")
	users.POST("/:username/follow", authMW, followHandler.Follow)
	users.DELETE("/:username/follow", authMW, followHandler.Unfollow)

	// feed
	r.GET("/feed", authMW, feedHandler.Get)

	protected := r.Group("/")
	protected.Use(middleware.Auth(cfg.JWTSecret))
	protected.GET("/me", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"user_id": c.GetUint(middleware.UserIDKey)})
	})

	log.Printf("server running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}
