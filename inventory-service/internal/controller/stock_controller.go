package controller

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/model"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/repository"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/usecases"
)

type StockController struct {
	stockUseCase usecase.StockUseCase
}

func NewStockController(usecase usecase.StockUseCase) StockController {
	return StockController{
		stockUseCase: usecase,
	}
}

func (sc *StockController) DeductStockController(ctx *gin.Context) {
	var deduction model.StockDeduction
	err := ctx.ShouldBindJSON(&deduction)
	if err != nil {
		response := model.Response{
			Message: "invalid request body",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	err = sc.stockUseCase.DeductStockUseCase(ctx.Request.Context(), deduction)
	if err != nil {
		status := http.StatusInternalServerError
		response := model.Response{
			Message: "could not deduct stock",
		}
		if errors.Is(err, usecase.ErrInvalidInvoiceId) ||
			errors.Is(err, usecase.ErrDeductionWithoutItems) ||
			errors.Is(err, usecase.ErrCodeRequired) ||
			errors.Is(err, usecase.ErrCodeTooLong) ||
			errors.Is(err, usecase.ErrInvalidQuantity) {
			status = http.StatusBadRequest
			response.Message = err.Error()
		} else if errors.Is(err, repository.ErrProductNotFound) {
			status = http.StatusNotFound
			response.Message = err.Error()
		} else if errors.Is(err, repository.ErrInsufficientBalance) ||
			errors.Is(err, repository.ErrDeductionConflict) {
			status = http.StatusConflict
			response.Message = err.Error()
		} else {
			log.Printf("stock deduction failed: %v", err)
		}
		ctx.JSON(status, response)
		return
	}
	response := model.Response{
		Message: "stock deducted successfully",
	}
	ctx.JSON(http.StatusOK, response)
}
