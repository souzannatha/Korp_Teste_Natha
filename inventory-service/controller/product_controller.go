package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/model"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/usecase"
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
	var product model.Product
	err := ctx.BindJSON(&product)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	insertedProduct, err := p.productUseCase.CreateProductUseCase(product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
	}

	ctx.JSON(http.StatusCreated, insertedProduct)
}

func (p *ProductController) GetProductController(ctx *gin.Context) {
	products, err := p.productUseCase.GetProductsUseCase()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, products)

}

func (p *ProductController) GetProductByCodeController(ctx *gin.Context) {

	code := ctx.Param("codeProduct")

	if code == "" {
		response := model.Response{
			Message: "Code not a nullable.",
		}
		ctx.JSON(http.StatusBadRequest, response)
		return
	}

	product, err := p.productUseCase.GetProductByCodeUseCase(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	if product == nil {
		response := model.Response{
			Message: "Product not found in data base.",
		}
		ctx.JSON(http.StatusNotFound, response)
		return
	}

	ctx.JSON(http.StatusOK, product)

}
