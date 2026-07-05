package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func getEnv(k, d string) string { if v := os.Getenv(k); v != "" { return v }; return d }

var (
	rateMu    sync.Mutex
	rateStore = map[string][]time.Time{}
)

func rateLimit(maxPerSec int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		now := time.Now()
		rateMu.Lock()
		times := rateStore[ip]
		var filtered []time.Time
		for _, t := range times {
			if now.Sub(t) < time.Second { filtered = append(filtered, t) }
		}
		filtered = append(filtered, now)
		rateStore[ip] = filtered
		rateMu.Unlock()
		if len(filtered) > maxPerSec {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
			return
		}
		c.Next()
	}
}

func jwtMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"}); return }
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid format"}); return }
		tok, err := jwt.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok { return nil, jwt.ErrSignatureInvalid }
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid { c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"}); return }
		cl := tok.Claims.(jwt.MapClaims)
		c.Set("user_id", cl["user_id"])
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" { c.AbortWithStatus(http.StatusNoContent); return }
		c.Next()
	}
}

func proxy(target string) gin.HandlerFunc {
	u, _ := url.Parse(target)
	p := httputil.NewSingleHostReverseProxy(u)
	return func(c *gin.Context) { p.ServeHTTP(c.Writer, c.Request) }
}

func main() {
	port       := getEnv("PORT", "8080")
	env        := getEnv("ENV", "development")
	jwtSecret  := getEnv("JWT_SECRET", "dev-jwt-secret-change-in-prod")
	userSvcURL := getEnv("USER_SERVICE_URL", "http://localhost:8081")
	prodSvcURL := getEnv("PRODUCT_SERVICE_URL", "http://localhost:8082")
	orderSvcURL := getEnv("ORDER_SERVICE_URL", "http://localhost:8083")
	paySvcURL  := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8084")

	logger, _ := zap.NewProduction()
	if env == "development" { logger, _ = zap.NewDevelopment() }
	defer logger.Sync()

	if env != "development" { gin.SetMode(gin.ReleaseMode) }
	r := gin.New()
	r.Use(gin.Recovery(), corsMiddleware(), rateLimit(100))

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "api-gateway"}) })

	// Public routes
	r.POST("/api/v1/auth/register", proxy(userSvcURL))
	r.POST("/api/v1/auth/login",    proxy(userSvcURL))
	r.GET("/api/v1/products",       proxy(prodSvcURL))
	r.GET("/api/v1/products/:id",   proxy(prodSvcURL))

	// Protected routes
	auth := r.Group("/api/v1")
	auth.Use(jwtMiddleware(jwtSecret))
	auth.GET("/users/me",      proxy(userSvcURL))
	auth.PUT("/users/me",      proxy(userSvcURL))
	auth.DELETE("/users/me",   proxy(userSvcURL))
	auth.POST("/orders",       proxy(orderSvcURL))
	auth.GET("/orders",        proxy(orderSvcURL))
	auth.GET("/orders/:id",    proxy(orderSvcURL))
	auth.DELETE("/orders/:id", proxy(orderSvcURL))
	auth.POST("/payments",     proxy(paySvcURL))
	auth.GET("/payments/:id",  proxy(paySvcURL))

	logger.Info("api-gateway starting", zap.String("port", port))
	srv := &http.Server{Addr: ":" + port, Handler: r, ReadTimeout: 30 * time.Second, WriteTimeout: 30 * time.Second}
	srv.ListenAndServe()
}
