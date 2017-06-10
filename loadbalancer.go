package main

import "gopkg.in/gin-gonic/gin.v1"
import "net/http/httputil"
import "net/url"
import "fmt"
import "math/rand"
import "time"

func main() {
	router := gin.Default()

	router.Any("/:path", ReverseProxy)

	router.Run(":8080")
}

func ReverseProxy(c *gin.Context) {
	target := RandomServer()

	url, err := url.Parse(target)
	fmt.Print(err)
	// checkErr(err)
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(c.Writer, c.Request)
}

func RandomServer() string {
	rand.Seed(time.Now().Unix()) // necesito algo mas random todavia

	target_list := []string{
		"http://localhost:8081",
		"http://localhost:8082",
	}
	n := rand.Int() % len(target_list)

	return target_list[n]
}
