package handler

import (
	"net/http"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func JWTMiddleware(secret string) gin.HandlerFunc {
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

func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		logger.Info("req", zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()), zap.Duration("latency", time.Since(start)))
	}
}

func MetricsHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) { h.ServeHTTP(c.Writer, c.Request) }
}
