package main

import "gopkg.in/gin-gonic/gin.v1"

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

    router.GET("/test", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "test server 2",
        })
    })

	router.Run(":8082")
}
