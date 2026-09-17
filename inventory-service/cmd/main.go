package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/controller"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/db"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/repository"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/routes"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/usecase"
)

func main() {
	server := gin.Default()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Não foi possível carregar o .env")
	}

	dbConnection, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}

	//Camada de repository
	ProductRepository := repository.NewProductRepository(dbConnection)

	//Camada usecase
	ProductUseCase := usecase.NewProductUseCase(ProductRepository)

	productController := controller.NewProductController(ProductUseCase)

	routes.RegisterProductRoutes(server, &productController)

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	server.Run(":8080")
}
