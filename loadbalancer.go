package main

import "gopkg.in/gin-gonic/gin.v1"
import "net/http/httputil"
import "net/url"
import "fmt"
import "math/rand"
import "time"

func main() {
	router := gin.Default()

    router.GET("/:route", ReverseProxy())
    router.POST("/:route", ReverseProxy())
    router.PUT("/:route", ReverseProxy())
    router.DELETE("/:route", ReverseProxy())

	router.GET("/test", ReverseProxy(RandomServer())) // ver como hacer para que la url sea wildcard

	router.Run(":8080") // listen and serve on 0.0.0.0:8080
}

func ReverseProxy(target string) gin.HandlerFunc {
    fmt.Print("-------------- exec reverse proxy --------------------\n")

    url, err := url.Parse(target)
    fmt.Print(err)
    // checkErr(err)
    proxy := httputil.NewSingleHostReverseProxy(url)
    return func(c *gin.Context) {
        proxy.ServeHTTP(c.Writer, c.Request)
    }
}

func RandomServer() string {
    fmt.Print("-------------- exec random server --------------------\n") // se esta ejecutando una unica vez

    rand.Seed(time.Now().Unix())

    target_list := []string{
        "http://localhost:8081",
        "http://localhost:8082",
    }
    n := rand.Int() % len(target_list)

    return target_list[n]
}
