package main
import ("context";"net/http";"os";"os/signal";"syscall";"time";"github.com/gin-gonic/gin";"github.com/yourname/ecommerce/order-service/internal/database";"github.com/yourname/ecommerce/order-service/internal/handler";"github.com/yourname/ecommerce/order-service/internal/messaging";"github.com/yourname/ecommerce/order-service/internal/repository";"github.com/yourname/ecommerce/order-service/internal/service";"go.uber.org/zap")
func getEnv(k,d string) string { if v:=os.Getenv(k);v!="" {return v}; return d }
func main() {
	port:=getEnv("PORT","8080"); env:=getEnv("ENV","development")
	dbURL:=getEnv("DATABASE_URL","postgres://order_svc:secret123@localhost:5434/orders_db?sslmode=disable")
	rabbitURL:=getEnv("RABBITMQ_URL","amqp://admin:admin123@localhost:5672/")
	logger,_:=zap.NewProduction(); if env=="development"{logger,_=zap.NewDevelopment()}; defer logger.Sync()
	db,err:=database.NewPostgres(dbURL); if err!=nil{logger.Fatal("db failed",zap.Error(err))}; defer db.Close()
	database.RunMigrations(dbURL)
	pub,err:=messaging.NewPublisher(rabbitURL,logger); if err!=nil{logger.Fatal("rabbitmq failed",zap.Error(err))}; defer pub.Close()
	repo:=repository.NewOrderRepository(db); svc:=service.NewOrderService(repo,pub,logger); h:=handler.NewOrderHandler(svc,logger)
	if env!="development"{gin.SetMode(gin.ReleaseMode)}
	r:=gin.New(); r.Use(gin.Recovery())
	r.GET("/health",func(c *gin.Context){c.JSON(http.StatusOK,gin.H{"status":"healthy","service":"order-service"})})
	v1:=r.Group("/api/v1"); v1.POST("/orders",h.CreateOrder); v1.GET("/orders",h.ListOrders); v1.GET("/orders/:id",h.GetOrder); v1.DELETE("/orders/:id",h.CancelOrder)
	srv:=&http.Server{Addr:":"+port,Handler:r,ReadTimeout:10*time.Second,WriteTimeout:10*time.Second}
	go func(){logger.Info("order-service starting",zap.String("port",port)); if err:=srv.ListenAndServe();err!=nil&&err!=http.ErrServerClosed{logger.Fatal("error",zap.Error(err))}}()
	quit:=make(chan os.Signal,1); signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM); <-quit
	ctx,cancel:=context.WithTimeout(context.Background(),30*time.Second); defer cancel(); srv.Shutdown(ctx)
}
