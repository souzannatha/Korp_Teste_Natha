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

type ProductController struct {
	productUseCase usecase.ProductUseCase
}

func NewProductController(usecase usecase.ProductUseCase) ProductController {
	return ProductController{
		productUseCase: usecase,
	}
}

func (p *ProductController) CreateProductController(ctx *gin.Context) {
	var request model.CreateProduct
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response := model.Response{
			Message: "invalid request body",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	if request.Balance == nil {
		response := model.Response{
			Message: usecase.ErrBalanceRequired.Error(),
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}
	product := model.Product{
		Code:        request.Code,
		Description: request.Description,
		Balance:     *request.Balance,
	}
	insertedProduct, err := p.productUseCase.CreateProductUseCase(ctx.Request.Context(), product)
	if err != nil {
		productErrorResponse(ctx, err, "could not create product")
		return
	}
	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *ProductController) GetProductController(ctx *gin.Context) {
	products, err := p.productUseCase.GetProductsUseCase(ctx.Request.Context())
	if err != nil {
		productErrorResponse(ctx, err, "could not retrieve products")
		return
	}
	ctx.JSON(http.StatusOK, products)
}

func (p *ProductController) GetProductByCodeController(ctx *gin.Context) {
	code := ctx.Param("codeProduct")
	product, err := p.productUseCase.GetProductByCodeUseCase(ctx.Request.Context(), code)
	if err != nil {
		productErrorResponse(ctx, err, "could not retrieve product")
		return
	}
	if product == nil {
		response := model.Response{
			Message: "product not found",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}
	ctx.JSON(http.StatusOK, product)
}

func productErrorResponse(ctx *gin.Context, err error, message string) {
	status := http.StatusInternalServerError
	response := model.Response{
		Message: message,
	}
	if errors.Is(err, usecase.ErrCodeRequired) ||
		errors.Is(err, usecase.ErrCodeTooLong) ||
		errors.Is(err, usecase.ErrDescriptionRequired) ||
		errors.Is(err, usecase.ErrInvalidBalance) {
		status = http.StatusBadRequest
		response.Message = err.Error()
	} else if errors.Is(err, repository.ErrDuplicateCode) {
		status = http.StatusConflict
		response.Message = err.Error()
	} else {
		log.Printf("product operation failed: %v", err)
	}
	ctx.JSON(status, response)
}
