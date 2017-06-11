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

func main() {
	config := LoadConfigFile("config.json")
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
	targetList := []string{
		"http://localhost:8081",
		"http://localhost:8082",
		"http://localhost:8083",
	}
	n := rand.Intn(100) % len(targetList)

	return targetList[n]
}

func LoadConfigFile(filename string) Config {
	var config Config
	configFile, _ := os.Open(filename)
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}
