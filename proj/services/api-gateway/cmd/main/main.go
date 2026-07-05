package main
import ("net/http";"net/http/httputil";"net/url";"os";"strings";"sync";"time";"github.com/gin-gonic/gin";"github.com/golang-jwt/jwt/v5";"go.uber.org/zap")
func getEnv(k,d string) string { if v:=os.Getenv(k);v!=""{return v};return d }
var (rateMu sync.Mutex; rateStore=map[string][]time.Time{})
func rateLimit(max int) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip:=c.ClientIP(); now:=time.Now(); rateMu.Lock()
		var f []time.Time; for _,t:=range rateStore[ip]{if now.Sub(t)<time.Second{f=append(f,t)}}
		f=append(f,now); rateStore[ip]=f; rateMu.Unlock()
		if len(f)>max{c.AbortWithStatusJSON(http.StatusTooManyRequests,gin.H{"error":"rate limit exceeded"});return}
		c.Next()
	}
}
func jwtMW(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h:=c.GetHeader("Authorization"); if h==""{c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"missing token"});return}
		parts:=strings.SplitN(h," ",2); if len(parts)!=2||parts[0]!="Bearer"{c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"invalid format"});return}
		tok,err:=jwt.Parse(parts[1],func(t *jwt.Token)(interface{},error){if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{return nil,jwt.ErrSignatureInvalid};return []byte(secret),nil})
		if err!=nil||!tok.Valid{c.AbortWithStatusJSON(http.StatusUnauthorized,gin.H{"error":"invalid token"});return}
		cl:=tok.Claims.(jwt.MapClaims); c.Set("user_id",cl["user_id"]); c.Next()
	}
}
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin","*"); c.Header("Access-Control-Allow-Methods","GET,POST,PUT,DELETE,OPTIONS"); c.Header("Access-Control-Allow-Headers","Authorization,Content-Type")
		if c.Request.Method=="OPTIONS"{c.AbortWithStatus(http.StatusNoContent);return}; c.Next()
	}
}
func proxy(target string) gin.HandlerFunc {
	u,_:=url.Parse(target); p:=httputil.NewSingleHostReverseProxy(u)
	return func(c *gin.Context){p.ServeHTTP(c.Writer,c.Request)}
}
func main() {
	port:=getEnv("PORT","8080"); env:=getEnv("ENV","development")
	secret:=getEnv("JWT_SECRET","dev-secret")
	userURL:=getEnv("USER_SERVICE_URL","http://localhost:8081"); prodURL:=getEnv("PRODUCT_SERVICE_URL","http://localhost:8082")
	orderURL:=getEnv("ORDER_SERVICE_URL","http://localhost:8083"); payURL:=getEnv("PAYMENT_SERVICE_URL","http://localhost:8084")
	logger,_:=zap.NewProduction(); if env=="development"{logger,_=zap.NewDevelopment()}; defer logger.Sync()
	if env!="development"{gin.SetMode(gin.ReleaseMode)}
	r:=gin.New(); r.Use(gin.Recovery(),cors(),rateLimit(100))
	r.GET("/health",func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"status":"healthy","service":"api-gateway"})})
	r.POST("/api/v1/auth/register",proxy(userURL)); r.POST("/api/v1/auth/login",proxy(userURL))
	r.GET("/api/v1/products",proxy(prodURL)); r.GET("/api/v1/products/:id",proxy(prodURL))
	auth:=r.Group("/api/v1"); auth.Use(jwtMW(secret))
	auth.GET("/users/me",proxy(userURL)); auth.PUT("/users/me",proxy(userURL)); auth.DELETE("/users/me",proxy(userURL))
	auth.POST("/orders",proxy(orderURL)); auth.GET("/orders",proxy(orderURL)); auth.GET("/orders/:id",proxy(orderURL)); auth.DELETE("/orders/:id",proxy(orderURL))
	auth.POST("/payments",proxy(payURL)); auth.GET("/payments/:id",proxy(payURL))
	logger.Info("api-gateway starting",zap.String("port",port))
	(&http.Server{Addr:":"+port,Handler:r,ReadTimeout:30*time.Second,WriteTimeout:30*time.Second}).ListenAndServe()
}
