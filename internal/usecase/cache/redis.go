package cache

import (
	"context"
	"encoding/json"
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
	Set(ctx context.Context, key string, quotation entity.CurrencyQuotation, ttl time.Duration) error
}

type cacheRedis struct {
	c *redis.Client
}

func NewCache(c *redis.Client) Cache {
	return &cacheRedis{c}
}

func (c *cacheRedis) Get(ctx context.Context, params GetParams) (*entity.CurrencyQuotation, error) {
	val := c.c.Get(ctx, params.CurrencyID+params.CurrencyIDTo+params.DateKey)

	if val.Err() != nil {
		return nil, val.Err()
	}

	var data []byte

	if err := val.Scan(&data); err != nil {
		return nil, err
	}

	out := new(entity.CurrencyQuotation)

	if err := json.Unmarshal(data, out); err != nil {
		return nil, err
	}

	return out, nil
}
func (c *cacheRedis) Set(
	ctx context.Context,
	key string,
	quotation entity.CurrencyQuotation,
	ttl time.Duration,
) error {
	data, err := quotation.MarshalBinary()
	if err != nil {
		return err
	}

	result := c.c.Set(ctx, key, data, ttl)

	return result.Err()
}
