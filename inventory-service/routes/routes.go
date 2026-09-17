package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/controller"
)

func RegisterProductRoutes(
	server *gin.Engine,
	pc *controller.ProductController,
) {
	products := server.Group("/product")

	products.POST("", pc.CreateProductController)
	products.GET("", pc.GetProductController)
	products.GET("/:codeProduct", pc.GetProductByCodeController)
}
