package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/ecommerce/payment-service/internal/domain"
)

type paymentRepository struct{ db *sqlx.DB }

func NewPaymentRepository(db *sqlx.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(p *domain.Payment) error {
	_, err := r.db.NamedExec(`INSERT INTO payments (id,order_id,stripe_intent_id,amount,currency,status,created_at,updated_at) VALUES (:id,:order_id,:stripe_intent_id,:amount,:currency,:status,:created_at,:updated_at)`, p)
	return err
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.Get(&p, `SELECT * FROM payments WHERE id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &p, err
}

func (r *paymentRepository) FindByOrderID(orderID string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.Get(&p, `SELECT * FROM payments WHERE order_id=$1`, orderID)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &p, err
}

func (r *paymentRepository) UpdateStatus(id uuid.UUID, status domain.PaymentStatus) error {
	_, err := r.db.Exec(`UPDATE payments SET status=$1,updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id)
	return err
}
