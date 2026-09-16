package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	server := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Não foi possível carregar o .env")
	}

	// dbConnection, err := db.ConnectDB()
	// if err != nil {
	// 	panic(err)
	// }

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.Run()
}
