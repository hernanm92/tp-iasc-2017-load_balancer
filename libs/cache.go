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

func (cacheClient CacheClient) SetRequest(request string, response string) {
	cacheClient.RedisClient.Set(request, response, 0)
}

func (cacheClient CacheClient) GetRequest(request string) (data string) {
	val, error := cacheClient.RedisClient.Get(request).Result()
	if error == redis.Nil {
		fmt.Println("Key no existe")
	}
	return string(val)
}

