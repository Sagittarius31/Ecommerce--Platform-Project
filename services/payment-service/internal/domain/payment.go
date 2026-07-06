package domain
import ("errors";"time";"github.com/google/uuid")
var ErrNotFound = errors.New("not found")
type PaymentStatus string
const (StatusPending PaymentStatus="pending"; StatusSucceeded PaymentStatus="succeeded"; StatusFailed PaymentStatus="failed")
type Payment struct {
	ID uuid.UUID `db:"id" json:"id"`; OrderID string `db:"order_id" json:"order_id"`
	StripeIntentID string `db:"stripe_intent_id" json:"stripe_intent_id"`
	Amount float64 `db:"amount" json:"amount"`; Currency string `db:"currency" json:"currency"`
	Status PaymentStatus `db:"status" json:"status"`; CreatedAt time.Time `db:"created_at" json:"created_at"`; UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type CreatePaymentInput struct { OrderID string `json:"order_id" validate:"required"`; Amount float64 `json:"amount" validate:"required,gt=0"` }
type PaymentRepository interface {
	Create(p *Payment) error; FindByID(id uuid.UUID) (*Payment, error)
	FindByOrderID(orderID string) (*Payment, error); UpdateStatus(id uuid.UUID, status PaymentStatus) error
}
