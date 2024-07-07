package cache

import (
	"advertisement-service/pkg/config"

	"github.com/go-redis/redis/v8"
)

func NewRedisCache(cfg config.CacheConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: cfg.Addr,
	})
}
