package controller

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/client"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/model"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/usecases"
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
	insertedInvoice, err := ic.invoiceUseCase.CreateInvoiceUseCase(ctx.Request.Context(), invoice)
	if err != nil {
		invoiceErrorResponse(ctx, err, "could not create invoice")
		return
	}
	ctx.JSON(http.StatusCreated, insertedInvoice)
}

func (ic *InvoiceController) GetAllInvoicesController(ctx *gin.Context) {
	invoices, err := ic.invoiceUseCase.GetAllInvoicesUseCase(ctx.Request.Context())
	if err != nil {
		invoiceErrorResponse(ctx, err, "could not retrieve invoices")
		return
	}
	ctx.JSON(http.StatusOK, invoices)
}

func (ic *InvoiceController) GetInvoiceByIdController(ctx *gin.Context) {
	invoiceId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := model.Response{
			Message: "invoice id must be a number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	invoice, err := ic.invoiceUseCase.GetInvoiceByIdUseCase(ctx.Request.Context(), invoiceId)
	if err != nil {
		invoiceErrorResponse(ctx, err, "could not retrieve invoice")
		return
	}
	ctx.JSON(http.StatusOK, invoice)
}

func (ic *InvoiceController) PrintInvoiceController(ctx *gin.Context) {
	invoiceId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response := model.Response{
			Message: "invoice id must be a number",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	invoice, err := ic.invoiceUseCase.PrintInvoiceUseCase(ctx.Request.Context(), invoiceId)
	if err != nil {
		invoiceErrorResponse(ctx, err, "could not print invoice")
		return
	}
	ctx.JSON(http.StatusOK, invoice)
}

func invoiceErrorResponse(ctx *gin.Context, err error, message string) {
	status := http.StatusInternalServerError
	response := model.Response{
		Message: message,
	}
	if errors.Is(err, usecase.ErrInvoiceWithoutItems) ||
		errors.Is(err, usecase.ErrProductCodeRequired) ||
		errors.Is(err, usecase.ErrProductCodeTooLong) ||
		errors.Is(err, usecase.ErrInvalidQuantity) ||
		errors.Is(err, usecase.ErrInvalidInvoiceId) ||
		errors.Is(err, client.ErrInventoryValidation) {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
		response.Message = "invoice not found"
	} else if errors.Is(err, client.ErrProductNotFound) {
		status = http.StatusNotFound
		response.Message = err.Error()
	} else if errors.Is(err, usecase.ErrInvoiceClosed) ||
		errors.Is(err, client.ErrInsufficientBalance) ||
		errors.Is(err, client.ErrDeductionConflict) {
		status = http.StatusConflict
		response.Message = err.Error()
	} else if errors.Is(err, client.ErrInventoryUnavailable) {
		status = http.StatusServiceUnavailable
		response.Message = err.Error()
	} else if errors.Is(err, client.ErrInventoryTimeout) {
		status = http.StatusGatewayTimeout
		response.Message = err.Error()
	} else if errors.Is(err, client.ErrInvalidInventoryResponse) {
		status = http.StatusBadGateway
		response.Message = err.Error()
	}
	if status >= 500 {
		log.Printf("invoice operation failed: %v", err)
	}
	ctx.JSON(status, response)
}
