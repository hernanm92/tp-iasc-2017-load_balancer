package main

import "gopkg.in/gin-gonic/gin.v1"
import "net/http/httputil"
import "net/url"
import "fmt"

func ReverseProxy() gin.HandlerFunc {
    target := "http://localhost:8081" // hacer que sea random entre 8081 y 8082

    url, err := url.Parse(target)
    fmt.Print(err)
    // checkErr(err)
    proxy := httputil.NewSingleHostReverseProxy(url)
    return func(c *gin.Context) {
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/test", ReverseProxy())

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}
