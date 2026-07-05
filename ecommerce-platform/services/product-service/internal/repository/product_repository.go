package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yourname/ecommerce/product-service/internal/domain"
)

type productRepository struct{ db *sqlx.DB }

func NewProductRepository(db *sqlx.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(p *domain.Product) error {
	_, err := r.db.NamedExec(`INSERT INTO products
		(id,name,description,price,stock,category_id,image_url,is_active,created_at,updated_at)
		VALUES (:id,:name,:description,:price,:stock,:category_id,:image_url,:is_active,:created_at,:updated_at)`, p)
	return err
}

func (r *productRepository) FindByID(id uuid.UUID) (*domain.Product, error) {
	var p domain.Product
	err := r.db.Get(&p, `SELECT * FROM products WHERE id=$1 AND is_active=true`, id)
	if errors.Is(err, sql.ErrNoRows) { return nil, domain.ErrNotFound }
	return &p, err
}

func (r *productRepository) List(f domain.ProductFilter) ([]*domain.Product, int, error) {
	conds := []string{"is_active=true"}
	args  := []interface{}{}
	i := 1
	if f.CategoryID != nil { conds = append(conds, fmt.Sprintf("category_id=$%d",i)); args = append(args,*f.CategoryID); i++ }
	if f.Search != ""      { conds = append(conds, fmt.Sprintf("name ILIKE $%d",i)); args = append(args,"%"+f.Search+"%"); i++ }
	where := "WHERE " + strings.Join(conds, " AND ")
	var total int
	r.db.QueryRow("SELECT COUNT(*) FROM products "+where, args...).Scan(&total)
	if f.Page == 0 { f.Page = 1 }
	if f.PageSize == 0 { f.PageSize = 20 }
	offset := (f.Page - 1) * f.PageSize
	rows, err := r.db.Queryx(fmt.Sprintf("SELECT * FROM products %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d",where,i,i+1), append(args,f.PageSize,offset)...)
	if err != nil { return nil, 0, err }
	defer rows.Close()
	var out []*domain.Product
	for rows.Next() { var p domain.Product; rows.StructScan(&p); out = append(out, &p) }
	return out, total, nil
}

func (r *productRepository) Update(p *domain.Product) error {
	_, err := r.db.NamedExec(`UPDATE products SET name=:name,description=:description,price=:price,stock=:stock,updated_at=:updated_at WHERE id=:id`, p)
	return err
}

func (r *productRepository) Delete(id uuid.UUID) error {
	_, err := r.db.Exec(`UPDATE products SET is_active=false WHERE id=$1`, id)
	return err
}

func (r *productRepository) DecrementStock(id uuid.UUID, qty int) error {
	res, err := r.db.Exec(`UPDATE products SET stock=stock-$1 WHERE id=$2 AND stock>=$1`, qty, id)
	if err != nil { return err }
	n, _ := res.RowsAffected()
	if n == 0 { return fmt.Errorf("insufficient stock") }
	return nil
}
