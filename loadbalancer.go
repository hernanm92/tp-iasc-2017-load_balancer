package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-gonic/gin"
)

type Config struct {
	Backends []string `json:"backends"`
}

var config Config

func main() {
	LoadConfigFile("config.json")
	fmt.Print(config)

	router := gin.Default()

	router.Any("/*path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(c *gin.Context) {
	target := RandomServer()

	url, err := url.Parse(target)
	fmt.Print(err)

	proxy := httputil.NewSingleHostReverseProxy(url)
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
