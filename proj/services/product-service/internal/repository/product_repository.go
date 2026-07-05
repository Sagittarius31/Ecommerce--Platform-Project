package repository
import ("database/sql"; "errors"; "fmt"; "github.com/google/uuid"; "github.com/jmoiron/sqlx"; "github.com/yourname/ecommerce/product-service/internal/domain")
type productRepo struct{ db *sqlx.DB }
func NewProductRepository(db *sqlx.DB) domain.ProductRepository { return &productRepo{db: db} }
func (r *productRepo) Create(p *domain.Product) error {
	_, err := r.db.NamedExec(`INSERT INTO products (id,name,description,price,stock,category_id,is_active,created_at,updated_at) VALUES (:id,:name,:description,:price,:stock,:category_id,:is_active,:created_at,:updated_at)`, p); return err
}
func (r *productRepo) FindByID(id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Get(&p, `SELECT * FROM products WHERE id=$1 AND is_active=true`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &p, err
}
func (r *productRepo) List(f domain.ProductFilter) ([]*domain.Product, int, error) {
	var total int; r.db.QueryRow(`SELECT COUNT(*) FROM products WHERE is_active=true`).Scan(&total)
	if f.Page == 0 { f.Page = 1 }; if f.PageSize == 0 { f.PageSize = 20 }
	var out []*domain.Product
	r.db.Select(&out, `SELECT * FROM products WHERE is_active=true ORDER BY created_at DESC LIMIT $1 OFFSET $2`, f.PageSize, (f.Page-1)*f.PageSize)
	return out, total, nil
}
func (r *productRepo) Update(p *domain.Product) error {
	_, err := r.db.NamedExec(`UPDATE products SET name=:name,description=:description,price=:price,stock=:stock,updated_at=:updated_at WHERE id=:id`, p); return err
}
func (r *productRepo) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE products SET is_active=false WHERE id=$1`, id); return err
}
func (r *productRepo) DecrementStock(id uuid.UUID, qty int) error {
	res, err := r.db.Exec(`UPDATE products SET stock=stock-$1 WHERE id=$2 AND stock>=$1`, qty, id)
	if err != nil { return err }
	n, _ := res.RowsAffected(); if n == 0 { return fmt.Errorf("insufficient stock") }; return nil
}
