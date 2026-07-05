package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourname/ecommerce/product-service/internal/cache"
	"github.com/yourname/ecommerce/product-service/internal/domain"
	"go.uber.org/zap"
)

type ProductService struct {
	repo   domain.ProductRepository
	cache  *cache.ProductCache
	logger *zap.Logger
}

func NewProductService(repo domain.ProductRepository, c *cache.ProductCache, logger *zap.Logger) *ProductService {
	return &ProductService{repo: repo, cache: c, logger: logger}
}

func (s *ProductService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if s.cache != nil {
		if cached := s.cache.Get(ctx, id.String()); cached != nil {
			s.logger.Debug("cache hit", zap.String("id", id.String()))
			return cached, nil
		}
	}
	p, err := s.repo.FindByID(id)
	if err != nil { return nil, err }
	if s.cache != nil { s.cache.Set(ctx, p) }
	return p, nil
}

func (s *ProductService) List(ctx context.Context, f domain.ProductFilter) ([]*domain.Product, int, error) {
	if f.Page == 0    { f.Page = 1 }
	if f.PageSize == 0 { f.PageSize = 20 }
	return s.repo.List(f)
}

func (s *ProductService) Create(ctx context.Context, in domain.CreateProductInput) (*domain.Product, error) {
	p := &domain.Product{
		ID: uuid.New(), Name: in.Name, Description: in.Description,
		Price: in.Price, Stock: in.Stock, CategoryID: in.CategoryID,
		ImageURL: in.ImageURL, IsActive: true,
		CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
	}
	if err := s.repo.Create(p); err != nil { return nil, err }
	if s.cache != nil { s.cache.Set(ctx, p) }
	return p, nil
}

func (s *ProductService) Update(ctx context.Context, id uuid.UUID, in domain.CreateProductInput) (*domain.Product, error) {
	p, err := s.repo.FindByID(id)
	if err != nil { return nil, err }
	p.Name = in.Name; p.Description = in.Description
	p.Price = in.Price; p.Stock = in.Stock
	p.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(p); err != nil { return nil, err }
	if s.cache != nil { s.cache.Invalidate(ctx, id.String()) }
	return p, nil
}

func (s *ProductService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.Delete(id); err != nil { return err }
	if s.cache != nil { s.cache.Invalidate(ctx, id.String()) }
	return nil
}
