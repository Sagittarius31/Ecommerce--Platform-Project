package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v78"
	"github.com/yourname/ecommerce/payment-service/internal/database"
	"github.com/yourname/ecommerce/payment-service/internal/handler"
	"github.com/yourname/ecommerce/payment-service/internal/repository"
	"github.com/yourname/ecommerce/payment-service/internal/service"
	"go.uber.org/zap"
)

func getEnv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

func main() {
	port          := getEnv("PORT", "8080")
	env           := getEnv("ENV", "development")
	dbURL         := getEnv("DATABASE_URL", "postgres://payment_svc:secret123@localhost:5435/payments_db?sslmode=disable")
	stripeKey     := getEnv("STRIPE_SECRET_KEY", "sk_test_REPLACE_ME")
	webhookSecret := getEnv("STRIPE_WEBHOOK_SECRET", "whsec_REPLACE_ME")
	stripe.Key = stripeKey

	logger, _ := zap.NewProduction()
	if env == "development" { logger, _ = zap.NewDevelopment() }
	defer logger.Sync()

	db, err := database.NewPostgres(dbURL)
	if err != nil { logger.Fatal("db failed", zap.Error(err)) }
	defer db.Close()
	database.RunMigrations(dbURL)

	repo := repository.NewPaymentRepository(db)
	svc  := service.NewPaymentService(repo, logger)
	h    := handler.NewPaymentHandler(svc, webhookSecret, logger)

	if env != "development" { gin.SetMode(gin.ReleaseMode) }
	r := gin.New()
	r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "payment-service"}) })
	r.POST("/webhooks/stripe", h.StripeWebhook)

	v1 := r.Group("/api/v1")
	v1.POST("/payments",      h.CreatePaymentIntent)
	v1.GET("/payments/:id",   h.GetPayment)

	srv := &http.Server{Addr: ":" + port, Handler: r, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}
	go func() {
		logger.Info("payment-service starting", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Fatal("server error", zap.Error(err)) }
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
