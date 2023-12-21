package cache

import (
	"github.com/redis/go-redis"
)

func NewRedisClient(host, port string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Network: host,
		Addr:    port,
	})

	return rdb
}
