package service
import ("context";"time";"github.com/google/uuid";"github.com/stripe/stripe-go/v78";"github.com/stripe/stripe-go/v78/paymentintent";"github.com/yourname/ecommerce/payment-service/internal/domain";"go.uber.org/zap")
type PaymentService struct { repo domain.PaymentRepository; logger *zap.Logger }
func NewPaymentService(repo domain.PaymentRepository, logger *zap.Logger) *PaymentService { return &PaymentService{repo:repo,logger:logger} }
func (s *PaymentService) CreatePaymentIntent(ctx context.Context, in domain.CreatePaymentInput) (*domain.Payment, *stripe.PaymentIntent, error) {
	params:=&stripe.PaymentIntentParams{Amount:stripe.Int64(int64(in.Amount*100)),Currency:stripe.String("usd"),Metadata:map[string]string{"order_id":in.OrderID}}
	intent,err:=paymentintent.New(params); if err!=nil{return nil,nil,err}
	p:=&domain.Payment{ID:uuid.New(),OrderID:in.OrderID,StripeIntentID:intent.ID,Amount:in.Amount,Currency:"usd",Status:domain.StatusPending,CreatedAt:time.Now().UTC(),UpdatedAt:time.Now().UTC()}
	if err:=s.repo.Create(p);err!=nil{return nil,nil,err}; return p,intent,nil
}
func (s *PaymentService) GetPayment(ctx context.Context, id uuid.UUID) (*domain.Payment, error) { return s.repo.FindByID(id) }
func (s *PaymentService) MarkSucceeded(orderID string) error {
	p,err:=s.repo.FindByOrderID(orderID); if err!=nil{return err}; return s.repo.UpdateStatus(p.ID,domain.StatusSucceeded)
}
func (s *PaymentService) MarkFailed(orderID string) error {
	p,err:=s.repo.FindByOrderID(orderID); if err!=nil{return err}; return s.repo.UpdateStatus(p.ID,domain.StatusFailed)
}
