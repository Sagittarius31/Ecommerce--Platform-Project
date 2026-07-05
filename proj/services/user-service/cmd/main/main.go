package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/yourname/ecommerce/user-service/internal/config"
	"github.com/yourname/ecommerce/user-service/internal/database"
	"github.com/yourname/ecommerce/user-service/internal/handler"
	"github.com/yourname/ecommerce/user-service/internal/repository"
	"github.com/yourname/ecommerce/user-service/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()
	logger, _ := zap.NewProduction()
	if cfg.Env == "development" { logger, _ = zap.NewDevelopment() }
	defer logger.Sync()

	db, err := database.NewPostgres(cfg.DatabaseURL)
	if err != nil { logger.Fatal("db failed", zap.Error(err)) }
	defer db.Close()
	database.RunMigrations(cfg.DatabaseURL)

	repo := repository.NewUserRepository(db)
	svc  := service.NewUserService(repo, cfg.JWTSecret, logger)
	h    := handler.NewUserHandler(svc, logger)

	if cfg.Env != "development" { gin.SetMode(gin.ReleaseMode) }
	r := gin.New()
	r.Use(gin.Recovery(), handler.LoggingMiddleware(logger))
	r.GET("/health",  func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"user-service"}) })
	r.GET("/metrics", handler.MetricsHandler())

	v1 := r.Group("/api/v1")
	v1.POST("/auth/register", h.Register)
	v1.POST("/auth/login",    h.Login)

	auth := v1.Group("/users")
	auth.Use(handler.JWTMiddleware(cfg.JWTSecret))
	auth.GET("/me",    h.GetProfile)
	auth.PUT("/me",    h.UpdateProfile)
	auth.DELETE("/me", h.DeleteAccount)

	srv := &http.Server{Addr: ":"+cfg.Port, Handler: r, ReadTimeout: 10*time.Second, WriteTimeout: 10*time.Second}
	go func() {
		logger.Info("user-service starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Fatal("listen error", zap.Error(err)) }
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
