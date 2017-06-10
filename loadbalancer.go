package main

import "gopkg.in/gin-gonic/gin.v1"
import "net/http/httputil"
import "net/url"
import "fmt"
import "math/rand"

func main() {
	router := gin.Default()

	router.Any("/:path", ReverseProxy)

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
	target_list := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}
	n := rand.Intn(100) % len(target_list)

	return target_list[n]
}
