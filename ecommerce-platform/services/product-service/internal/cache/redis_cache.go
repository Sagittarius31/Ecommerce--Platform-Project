package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yourname/ecommerce/product-service/internal/domain"
	"go.uber.org/zap"
)

type ProductCache struct {
	client *redis.Client
	logger *zap.Logger
}

func NewProductCache(redisURL string, logger *zap.Logger) (*ProductCache, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil { return nil, err }
	client := redis.NewClient(opts)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil { return nil, fmt.Errorf("redis ping: %w", err) }
	return &ProductCache{client: client, logger: logger}, nil
}

func (c *ProductCache) Get(ctx context.Context, id string) *domain.Product {
	data, err := c.client.Get(ctx, "product:"+id).Bytes()
	if err != nil { return nil }
	var p domain.Product
	if err := json.Unmarshal(data, &p); err != nil { return nil }
	return &p
}

func (c *ProductCache) Set(ctx context.Context, p *domain.Product) {
	data, _ := json.Marshal(p)
	c.client.Set(ctx, "product:"+p.ID.String(), data, 10*time.Minute)
}

func (c *ProductCache) Invalidate(ctx context.Context, id string) {
	c.client.Del(ctx, "product:"+id)
}
