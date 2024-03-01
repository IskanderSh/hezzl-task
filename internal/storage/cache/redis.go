package redis

import (
	"context"
	"fmt"

	"github.com/IskanderSh/hezzl-task/internal/config"
	"github.com/IskanderSh/hezzl-task/internal/lib/error/wrapper"
	"github.com/IskanderSh/hezzl-task/internal/models"
	"github.com/redis/go-redis/v9"
)

const (
	priorityKey      = "priority"
	zeroExpiration   = 0
	minuteExpiration = 60
)

type Cache struct {
	client *redis.Client
}

func NewCache(cfg config.Cache) *Cache {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
	})

	return &Cache{client: client}
}

func (c *Cache) GetMaxPriority(ctx context.Context) (string, error) {
	const op = "storage.cache.GetMaxPriority"

	value, err := c.client.Get(ctx, priorityKey).Result()
	if err != nil {
		return "", wrapper.Wrap(op, err)
	}

	return value, nil
}

func (c *Cache) SetMaxPriority(ctx context.Context, priority int) error {
	const op = "storage.cache.SetMaxPriority"

	if err := c.client.Set(ctx, priorityKey, priority, zeroExpiration).Err(); err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}

func (c *Cache) SaveGood(ctx context.Context, key string, value *models.GoodCache) error {
	const op = "storage.cache.SaveGoods"

	if err := c.client.Set(ctx, key, value, minuteExpiration).Err(); err != nil {
		return wrapper.Wrap(op, err)
	}

	return nil
}

func (c *Cache) GetGood(ctx context.Context, key string) (string, error) {
	const op = "storage.cache.GetGood"

	value, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return "", wrapper.Wrap(op, err)
	}

	return value, nil
}
