package handler
import ("io";"net/http";"github.com/gin-gonic/gin";"github.com/google/uuid";"github.com/stripe/stripe-go/v78/webhook";"github.com/yourname/ecommerce/payment-service/internal/domain";"github.com/yourname/ecommerce/payment-service/internal/service";"go.uber.org/zap")
type PaymentHandler struct { svc *service.PaymentService; webhookSecret string; logger *zap.Logger }
func NewPaymentHandler(svc *service.PaymentService, secret string, logger *zap.Logger) *PaymentHandler { return &PaymentHandler{svc:svc,webhookSecret:secret,logger:logger} }
func (h *PaymentHandler) CreatePaymentIntent(c *gin.Context) {
	var in domain.CreatePaymentInput; if err:=c.ShouldBindJSON(&in);err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid body"});return}
	p,intent,err:=h.svc.CreatePaymentIntent(c.Request.Context(),in); if err!=nil{c.JSON(http.StatusInternalServerError,gin.H{"error":err.Error()});return}
	c.JSON(http.StatusCreated,gin.H{"payment":p,"client_secret":intent.ClientSecret})
}
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id,err:=uuid.Parse(c.Param("id")); if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid id"});return}
	p,err:=h.svc.GetPayment(c.Request.Context(),id); if err!=nil{c.JSON(http.StatusNotFound,gin.H{"error":"not found"});return}
	c.JSON(http.StatusOK,p)
}
func (h *PaymentHandler) StripeWebhook(c *gin.Context) {
	payload,_:=io.ReadAll(c.Request.Body)
	event,err:=webhook.ConstructEvent(payload,c.GetHeader("Stripe-Signature"),h.webhookSecret)
	if err!=nil{c.JSON(http.StatusBadRequest,gin.H{"error":"invalid signature"});return}
	switch event.Type {
	case "payment_intent.succeeded": h.logger.Info("payment succeeded")
	case "payment_intent.payment_failed": h.logger.Info("payment failed")
	}
	c.Status(http.StatusOK)
}
