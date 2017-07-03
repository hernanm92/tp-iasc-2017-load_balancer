package cache

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

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

func (cacheClient CacheClient) SetRequest(request *http.Request, responseData string, expiredTime int) {
	key := CreateRequesString(request)
	cacheClient.RedisClient.Set(key, responseData, time.Duration(expiredTime)*time.Minute)
}

func (cacheClient CacheClient) GetRequestValue(request *http.Request) (data string) {
	key := CreateRequesString(request)
	val, error := cacheClient.RedisClient.Get(key).Result()
	if error == redis.Nil {
		fmt.Println("Key no existe")
	}
	return string(val)
}

func (cacheClient CacheClient) ExistsOrNotExpiredKey(request *http.Request) bool {
	key := CreateRequesString(request)
	_, error := cacheClient.RedisClient.Get(key).Result()
	if error == redis.Nil {
		return false
	}
	return true
}

func (cacheClient CacheClient) IsCacheble(request *http.Request) bool {
	cachecontrol := string(request.Header.Get("Cache-Control"))
	fmt.Println(cachecontrol)
	return strings.EqualFold(string(request.Method), "GET") && !strings.EqualFold(cachecontrol, "no-cache") && !strings.EqualFold(cachecontrol, "expired")
}

func CreateRequesString(request *http.Request) (request_string string) {
	req, _ := httputil.DumpRequest(request, true)
	return string(req)
}
