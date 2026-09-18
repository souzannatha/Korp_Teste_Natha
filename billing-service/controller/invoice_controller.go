package controller

import (
	"errors"
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

	err := ctx.ShouldBindJSON(&invoice)

	if err != nil {
		response := model.Response{
			Message: "invalid request body",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	insertedInvoice, err := ic.invoiceUseCase.CreateInvoiceUseCase(invoice)

	if err != nil {
		if errors.Is(err, usecase.ErrInvoiceWithoutItems) ||
			errors.Is(err, usecase.ErrProductCodeRequired) ||
			errors.Is(err, usecase.ErrInvalidQuantity) {
			response := model.Response{
				Message: err.Error(),
			}
			ctx.JSON(http.StatusBadRequest, response)
			return
		}

		response := model.Response{
			Message: "could not create invoice",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}

	ctx.JSON(http.StatusCreated, insertedInvoice)
}

func (ic *InvoiceController) GetAllInvoicesController(ctx *gin.Context) {
	invoices, err := ic.invoiceUseCase.GetAllInvoicesUseCase()
	if err != nil {
		response := model.Response{
			Message: "could not retrieve invoices",
		}
		ctx.JSON(http.StatusInternalServerError, response)
		return
	}
	ctx.JSON(http.StatusOK, invoices)
}
