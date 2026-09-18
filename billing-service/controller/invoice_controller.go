package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/model"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/usecase"
)

type InvoiceController struct {
	invoiceUseCase usecase.InvoiceUseCase
}

func NewInvoiceController(usecase usecase.InvoiceUseCase) InvoiceController {
	return InvoiceController{
		invoiceUseCase: usecase,
	}
}

func (ic *InvoiceController) CreateInvoiceController(ctx *gin.Context) {
	var invoice model.Invoice

	err := ctx.BindJSON(&invoice)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedInvoice, err := ic.invoiceUseCase.CreateInvoiceUseCase(invoice)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, insertedInvoice)
}
