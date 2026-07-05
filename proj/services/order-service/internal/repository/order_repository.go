package repository
import ("database/sql";"errors";"time";"github.com/google/uuid";"github.com/jmoiron/sqlx";"github.com/yourname/ecommerce/order-service/internal/domain")
type orderRepo struct{ db *sqlx.DB }
func NewOrderRepository(db *sqlx.DB) domain.OrderRepository { return &orderRepo{db:db} }
func (r *orderRepo) Create(o *domain.Order, items []domain.OrderItem) error {
	tx, err := r.db.Beginx(); if err != nil { return err }; defer tx.Rollback()
	if _, err := tx.NamedExec(`INSERT INTO orders (id,user_id,user_email,status,total,created_at,updated_at) VALUES (:id,:user_id,:user_email,:status,:total,:created_at,:updated_at)`, o); err != nil { return err }
	for _, item := range items { if _, err := tx.NamedExec(`INSERT INTO order_items (id,order_id,product_id,quantity,price) VALUES (:id,:order_id,:product_id,:quantity,:price)`, item); err != nil { return err } }
	return tx.Commit()
}
func (r *orderRepo) FindByID(id uuid.UUID) (*domain.Order, error) {
	var o domain.Order; err := r.db.Get(&o, `SELECT * FROM orders WHERE id=$1`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }; return &o, err
}
func (r *orderRepo) UpdateStatus(id uuid.UUID, status domain.OrderStatus) error {
	_, err := r.db.Exec(`UPDATE orders SET status=$1,updated_at=$2 WHERE id=$3`, status, time.Now().UTC(), id); return err
}
