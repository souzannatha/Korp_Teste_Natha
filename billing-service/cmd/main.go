package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	server := gin.Default()

	server.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": true,
		})
	})

	server.Run(":8081")
}
