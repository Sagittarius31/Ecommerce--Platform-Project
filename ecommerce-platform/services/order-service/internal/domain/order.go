package domain

import (
	"errors"
	"time"
	"github.com/google/uuid"
)

var (
	ErrNotFound        = errors.New("not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusConfirmed OrderStatus = "confirmed"
	StatusCancelled OrderStatus = "cancelled"
	StatusFailed    OrderStatus = "payment_failed"
)

type Order struct {
	ID         uuid.UUID   `db:"id"         json:"id"`
	UserID     uuid.UUID   `db:"user_id"    json:"user_id"`
	Status     OrderStatus `db:"status"     json:"status"`
	Total      float64     `db:"total"      json:"total"`
	UserEmail  string      `db:"user_email" json:"user_email"`
	CreatedAt  time.Time   `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time   `db:"updated_at" json:"updated_at"`
	Items      []OrderItem `db:"-"          json:"items,omitempty"`
}

type OrderItem struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	OrderID   uuid.UUID `db:"order_id"   json:"order_id"`
	ProductID uuid.UUID `db:"product_id" json:"product_id"`
	Quantity  int       `db:"quantity"   json:"quantity"`
	Price     float64   `db:"price"      json:"price"`
}

type CreateOrderInput struct {
	UserID    uuid.UUID         `json:"user_id"    validate:"required"`
	UserEmail string            `json:"user_email" validate:"required,email"`
	Items     []CreateOrderItem `json:"items"      validate:"required,min=1"`
}

type CreateOrderItem struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity"   validate:"required,min=1"`
	Price     float64   `json:"price"      validate:"required,gt=0"`
}

type OrderRepository interface {
	Create(o *Order, items []OrderItem) error
	FindByID(id uuid.UUID) (*Order, error)
	ListByUser(userID uuid.UUID, page, pageSize int) ([]*Order, int, error)
	UpdateStatus(id uuid.UUID, status OrderStatus) error
}
