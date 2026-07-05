package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/ecommerce/order-service/internal/domain"
)

type orderRepository struct{ db *sqlx.DB }

func NewOrderRepository(db *sqlx.DB) domain.OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(o *domain.Order, items []domain.OrderItem) error {
	tx, err := r.db.Beginx()
	if err != nil { return err }
	defer tx.Rollback()
	if _, err := tx.NamedExec(`INSERT INTO orders (id,user_id,status,total,user_email,created_at,updated_at) VALUES (:id,:user_id,:status,:total,:user_email,:created_at,:updated_at)`, o); err != nil { return err }
	for _, item := range items {
		if _, err := tx.NamedExec(`INSERT INTO order_items (id,order_id,product_id,quantity,price) VALUES (:id,:order_id,:product_id,:quantity,:price)`, item); err != nil { return err }
	}
	return tx.Commit()
}

func (r *orderRepository) FindByID(id uuid.UUID) (*domain.Order, error) {
	var o domain.Order
	err := r.db.Get(&o, `SELECT * FROM orders WHERE id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &o, err
}

func (r *orderRepository) ListByUser(userID uuid.UUID, page, pageSize int) ([]*domain.Order, int, error) {
	var total int
	r.db.QueryRow(`SELECT COUNT(*) FROM orders WHERE user_id=$1`, userID).Scan(&total)
	offset := (page - 1) * pageSize
	var orders []*domain.Order
	err := r.db.Select(&orders, `SELECT * FROM orders WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`, userID, pageSize, offset)
	return orders, total, err
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status domain.OrderStatus) error {
	_, err := r.db.Exec(`UPDATE orders SET status=$1,updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id)
	return err
}
