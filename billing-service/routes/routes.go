package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/controller"
)

func RegisterInvoiceRoutes(
	server *gin.Engine,
	ic *controller.InvoiceController,
) {
	invoices := server.Group("/invoice")

	invoices.POST("", ic.CreateInvoiceController)
	invoices.GET("", ic.GetAllInvoicesController)
}
