package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/controller"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/db"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/repository"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/routes"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/usecase"
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
	invoiceRepository := repository.NewInvoiceRepository(dbConnection)

	// Camada de usecase
	invoiceUseCase := usecase.NewInvoiceUseCase(invoiceRepository)

	// Camada de controller
	invoiceController := controller.NewInvoiceController(invoiceUseCase)

	routes.RegisterInvoiceRoutes(server, &invoiceController)

	server.GET("/healthy", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health": true,
		})
	})

	server.Run(":8081")
}
