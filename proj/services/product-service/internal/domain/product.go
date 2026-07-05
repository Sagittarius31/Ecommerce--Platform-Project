package domain
import ("errors"; "time"; "github.com/google/uuid")
var ErrNotFound = errors.New("not found")
type Product struct {
	ID          uuid.UUID `db:"id"          json:"id"`
	Name        string    `db:"name"        json:"name"`
	Description string    `db:"description" json:"description"`
	Price       float64   `db:"price"       json:"price"`
	Stock       int       `db:"stock"       json:"stock"`
	CategoryID  uuid.UUID `db:"category_id" json:"category_id"`
	IsActive    bool      `db:"is_active"   json:"is_active"`
	CreatedAt   time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"  json:"updated_at"`
}
type CreateProductInput struct {
	Name        string    `json:"name"        validate:"required"`
	Description string    `json:"description" validate:"required"`
	Price       float64   `json:"price"       validate:"required,gt=0"`
	Stock       int       `json:"stock"       validate:"gte=0"`
	CategoryID  uuid.UUID `json:"category_id"`
}
type ProductFilter struct { Search string; Page, PageSize int }
type ProductRepository interface {
	Create(p *Product) error
	FindByID(id uuid.UUID) (*Product, error)
	List(f ProductFilter) ([]*Product, int, error)
	Update(p *Product) error
	Delete(id uuid.UUID) error
	DecrementStock(id uuid.UUID, qty int) error
}
