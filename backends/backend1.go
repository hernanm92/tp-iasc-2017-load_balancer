package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"time"
)

func main() {
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		time.Sleep(10 * time.Second)

		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/test", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "backend 1",
		})
	})

	router.Run(":8081")
}
