package main
import ("context"; "net/http"; "os"; "os/signal"; "syscall"; "time"; "github.com/gin-gonic/gin"; "github.com/yourname/ecommerce/product-service/internal/cache"; "github.com/yourname/ecommerce/product-service/internal/database"; "github.com/yourname/ecommerce/product-service/internal/handler"; "github.com/yourname/ecommerce/product-service/internal/repository"; "github.com/yourname/ecommerce/product-service/internal/service"; "go.uber.org/zap")
func getEnv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }
func main() {
	port := getEnv("PORT","8080"); env := getEnv("ENV","development")
	dbURL := getEnv("DATABASE_URL","postgres://product_svc:secret123@localhost:5433/products_db?sslmode=disable")
	redisURL := getEnv("REDIS_URL","redis://localhost:6379")
	logger, _ := zap.NewProduction(); if env == "development" { logger, _ = zap.NewDevelopment() }; defer logger.Sync()
	db, err := database.NewPostgres(dbURL); if err != nil { logger.Fatal("db failed", zap.Error(err)) }; defer db.Close()
	database.RunMigrations(dbURL)
	productCache, err := cache.NewProductCache(redisURL, logger); if err != nil { logger.Warn("redis unavailable", zap.Error(err)) }
	repo := repository.NewProductRepository(db); svc := service.NewProductService(repo, productCache, logger); h := handler.NewProductHandler(svc, logger)
	if env != "development" { gin.SetMode(gin.ReleaseMode) }
	r := gin.New(); r.Use(gin.Recovery())
	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status":"healthy","service":"product-service"}) })
	v1 := r.Group("/api/v1")
	v1.GET("/products", h.ListProducts); v1.GET("/products/:id", h.GetProduct)
	v1.POST("/products", h.CreateProduct); v1.PUT("/products/:id", h.UpdateProduct); v1.DELETE("/products/:id", h.DeleteProduct)
	srv := &http.Server{Addr:":"+port, Handler:r, ReadTimeout:10*time.Second, WriteTimeout:10*time.Second}
	go func() { logger.Info("product-service starting", zap.String("port",port)); if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed { logger.Fatal("error",zap.Error(err)) } }()
	quit := make(chan os.Signal,1); signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM); <-quit
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second); defer cancel(); srv.Shutdown(ctx)
}
