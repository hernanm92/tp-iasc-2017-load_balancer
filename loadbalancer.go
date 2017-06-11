package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httputil"
	"net/url"
	"os"

	"./libs"
	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends []string `json:"backends"`
	WaitTime int      `json:"wait_time"`
}

var config Config

func main() {
	LoadConfigFile("config.json")
	fmt.Println(config)

	fmt.Println("Starting Server...")

	router := gin.Default()

	router.Any("/*path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(c *gin.Context) {
	target := RandomServer()

	url, _ := url.Parse(target)

	proxy := httputil.NewSingleHostReverseProxy(url)
	// voy a necesitar saber la respuesta para poder cachearla

	cache_client := cache.CacheClient{cache.NewClient()}
	cache_client.SetRequest()
	proxy.ServeHTTP(c.Writer, c.Request)
}

func RandomServer() string {
	n := rand.Intn(100) % len(config.Backends)

	return config.Backends[n]
}

func LoadConfigFile(filename string) {
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
}
