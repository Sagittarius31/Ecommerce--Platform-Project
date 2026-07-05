package service
import ("context";"time";"github.com/google/uuid";"github.com/yourname/ecommerce/order-service/internal/domain";"github.com/yourname/ecommerce/order-service/internal/messaging";"go.uber.org/zap")
type OrderPlacedEvent struct { OrderID,UserID,UserEmail string; TotalAmount float64; Timestamp string }
type OrderService struct { repo domain.OrderRepository; pub *messaging.Publisher; logger *zap.Logger }
func NewOrderService(repo domain.OrderRepository, pub *messaging.Publisher, logger *zap.Logger) *OrderService {
	return &OrderService{repo:repo,pub:pub,logger:logger}
}
func (s *OrderService) CreateOrder(ctx context.Context, in domain.CreateOrderInput) (*domain.Order, error) {
	var total float64; var items []domain.OrderItem
	for _, i := range in.Items { total += i.Price*float64(i.Quantity); items = append(items, domain.OrderItem{ID:uuid.New(),ProductID:i.ProductID,Quantity:i.Quantity,Price:i.Price}) }
	o := &domain.Order{ID:uuid.New(),UserID:in.UserID,UserEmail:in.UserEmail,Status:domain.StatusPending,Total:total,CreatedAt:time.Now().UTC(),UpdatedAt:time.Now().UTC()}
	for i := range items { items[i].OrderID = o.ID }
	if err := s.repo.Create(o, items); err != nil { return nil, err }
	s.pub.Publish(ctx,"order.placed",OrderPlacedEvent{OrderID:o.ID.String(),UserID:o.UserID.String(),UserEmail:o.UserEmail,TotalAmount:o.Total,Timestamp:time.Now().UTC().Format(time.RFC3339)})
	s.logger.Info("order created", zap.String("id", o.ID.String())); return o, nil
}
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*domain.Order, error) { return s.repo.FindByID(id) }
func (s *OrderService) CancelOrder(ctx context.Context, id uuid.UUID) error { return s.repo.UpdateStatus(id, domain.StatusCancelled) }
