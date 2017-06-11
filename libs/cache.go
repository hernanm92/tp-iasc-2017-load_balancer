package cache

import (
	"fmt"

	"github.com/go-redis/redis"
)

type CacheClient struct {
	RedisClient *redis.Client
}

func NewClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return client
}

func (cacheClient CacheClient) SetRequest() {
	cacheClient.RedisClient.Set("key", "value", 0)
}

func (cacheClient CacheClient) GetRequest() {
	fmt.Println("exec cache.GetRequest")
}
