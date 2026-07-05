package main
import ("context";"net/http";"os";"os/signal";"syscall";"time";"github.com/gin-gonic/gin";"github.com/yourname/ecommerce/notification-service/internal/consumer";"github.com/yourname/ecommerce/notification-service/internal/email";"go.uber.org/zap")
func getEnv(k,d string) string { if v:=os.Getenv(k);v!=""{return v};return d }
func main() {
	port:=getEnv("PORT","8080"); env:=getEnv("ENV","development")
	rabbitURL:=getEnv("RABBITMQ_URL","amqp://admin:admin123@localhost:5672/")
	logger,_:=zap.NewProduction(); if env=="development"{logger,_=zap.NewDevelopment()}; defer logger.Sync()
	emailSvc:=email.New(email.Config{Host:getEnv("SMTP_HOST","localhost"),Port:getEnv("SMTP_PORT","587"),Username:getEnv("SMTP_USER",""),Password:getEnv("SMTP_PASSWORD",""),From:getEnv("SMTP_FROM","noreply@yourdomain.com")},logger)
	ctx,cancel:=context.WithCancel(context.Background()); defer cancel()
	nc,err:=consumer.New(rabbitURL,emailSvc,logger); if err!=nil{logger.Fatal("rabbitmq failed",zap.Error(err))}; defer nc.Close()
	if err:=nc.Start(ctx);err!=nil{logger.Fatal("consumer failed",zap.Error(err))}
	if env!="development"{gin.SetMode(gin.ReleaseMode)}
	r:=gin.New(); r.Use(gin.Recovery())
	r.GET("/health",func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"status":"healthy","service":"notification-service"})})
	srv:=&http.Server{Addr:":"+port,Handler:r,ReadTimeout:10*time.Second,WriteTimeout:10*time.Second}
	go func(){logger.Info("notification-service starting",zap.String("port",port)); if err:=srv.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{logger.Fatal("error",zap.Error(err))}}()
	quit:=make(chan os.Signal,1); signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM); <-quit
	cancel(); shutCtx,shutCancel:=context.WithTimeout(context.Background(),30*time.Second); defer shutCancel(); srv.Shutdown(shutCtx)
}
