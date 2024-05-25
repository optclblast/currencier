package cache

import (
	"context"
	"time"

	"github.com/optclblast/currencier/internal/entity"
	"github.com/redis/go-redis/v9"
)

type GetParams struct {
	CurrencyID   string
	CurrencyIDTo string
	DateKey      string
}

type Cache interface {
	Get(ctx context.Context, params GetParams) (*entity.CurrencyQuotation, error)
	Set(ctx context.Context, quotation entity.CurrencyQuotation, ttl time.Duration) error
}

type cacheRedis struct {
	c *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &cacheRedis{c}
}

func (c *cacheRedis) Get(ctx context.Context, params GetParams) (*entity.CurrencyQuotation, error) {
	return nil, nil
}
func (c *cacheRedis) Set(
	ctx context.Context,
	quotation entity.CurrencyQuotation,
	ttl time.Duration,
) error {
	return nil
}
