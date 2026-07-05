package domain
import ("errors";"time";"github.com/google/uuid")
var ErrNotFound = errors.New("not found")
type OrderStatus string
const (StatusPending OrderStatus="pending"; StatusConfirmed OrderStatus="confirmed"; StatusCancelled OrderStatus="cancelled")
type Order struct {
	ID uuid.UUID `db:"id" json:"id"`; UserID uuid.UUID `db:"user_id" json:"user_id"`
	UserEmail string `db:"user_email" json:"user_email"`; Status OrderStatus `db:"status" json:"status"`
	Total float64 `db:"total" json:"total"`; CreatedAt time.Time `db:"created_at" json:"created_at"`; UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
type OrderItem struct {
	ID uuid.UUID `db:"id" json:"id"`; OrderID uuid.UUID `db:"order_id" json:"order_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`; Quantity int `db:"quantity" json:"quantity"`; Price float64 `db:"price" json:"price"`
}
type CreateOrderInput struct { UserID uuid.UUID `json:"user_id"`; UserEmail string `json:"user_email"`; Items []CreateOrderItem `json:"items"` }
type CreateOrderItem struct { ProductID uuid.UUID `json:"product_id"`; Quantity int `json:"quantity"`; Price float64 `json:"price"` }
type OrderRepository interface {
	Create(o *Order, items []OrderItem) error
	FindByID(id uuid.UUID) (*Order, error)
	UpdateStatus(id uuid.UUID, status OrderStatus) error
}
