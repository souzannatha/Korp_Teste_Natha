package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/controller"
)

func RegisterStockRoutes(server *gin.Engine, sc *controller.StockController) {
	server.POST("/stock-deductions", sc.DeductStockController)
}
