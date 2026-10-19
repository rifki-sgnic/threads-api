package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"

	"github.com/gin-gonic/gin"
	"github.com/rifki-sgnic/threads-api/internal/config"
	"github.com/rifki-sgnic/threads-api/internal/database"
	"github.com/rifki-sgnic/threads-api/internal/docs"
	"github.com/rifki-sgnic/threads-api/internal/handler"
	"github.com/rifki-sgnic/threads-api/internal/middleware"
	"github.com/rifki-sgnic/threads-api/internal/repository"
	"github.com/rifki-sgnic/threads-api/internal/service"
)

func main() {
	logger := newLogger()
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("invalid config", "error", err)
		os.Exit(1)
	}

	db, err := database.Connect(cfg)
	if err != nil {
		logger.Error("failed to connect database", "error", err)
		os.Exit(1)
	}

	if err := database.Migrate(db); err != nil {
		logger.Error("failed to migrate database: %v", err)
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
	profileSvc := service.NewProfileService(userRepo, followRepo, postRepo, statsLoader)

	// handlers
	authHandler := handler.NewAuthHandler(authSvc)
	postHandler := handler.NewPostHandler(postSvc)
	likeHandler := handler.NewLikeHandler(likeSvc)
	followHandler := handler.NewFollowHandler(followSvc)
	feedHandler := handler.NewFeedHandler(feedSvc)
	profileHandler := handler.NewProfileHandler(profileSvc)

	authMW := middleware.Auth(cfg.JWTSecret)
	optAuthMW := middleware.OptionalAuth(cfg.JWTSecret)

	r := gin.New()
	r.Use(
		cors.New(cors.Config{
			AllowOrigins:  cfg.CORSOrigins,
			AllowMethods:  []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:  []string{"Authorization", "Content-Type"},
			ExposeHeaders: []string{"X-Request-ID"},
			MaxAge:        12 * time.Hour,
		}),
		middleware.RequestLogger(logger),
		middleware.Recovery(),
	)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unavailable"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	if cfg.EnableDocs {
		docs.Register(r)
	}

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
	users.GET("/:username", optAuthMW, profileHandler.Get)
	users.GET("/:username/posts", optAuthMW, profileHandler.Posts)
	users.POST("/:username/follow", authMW, followHandler.Follow)
	users.DELETE("/:username/follow", authMW, followHandler.Unfollow)

	r.GET("/me", authMW, profileHandler.Me)
	r.PATCH("/me", authMW, profileHandler.UpdateMe)

	// feed
	r.GET("/feed", authMW, feedHandler.Get)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		logger.Info("server running", "port", cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	logger.Info("shutting down, waiting for in-flight requests")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("forced shutdown", "error", err)
	}

	if sqlDB, err := db.DB(); err == nil {
		_ = sqlDB.Close()
	}
	logger.Info("bye")
}

func newLogger() *slog.Logger {
	if os.Getenv("GIN_MODE") == "release" {
		return slog.New(slog.NewJSONHandler(os.Stdout, nil)) // production: JSON
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil)) // dev: enak dibaca
}
