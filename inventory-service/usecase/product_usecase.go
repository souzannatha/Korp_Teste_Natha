package usecase

import (
	"errors"
	"strings"

	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/model"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/repository"
)

type ProductUseCase struct {
	repository repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return ProductUseCase{
		repository: repo,
	}
}

func (pu *ProductUseCase) CreateProductUseCase(product model.Product) (model.Product, error) {
	product.Code = strings.TrimSpace(product.Code)
	product.Description = strings.TrimSpace(product.Description)

	if product.Code == "" {
		return model.Product{}, errors.New("code is required")
	}

	if product.Description == "" {
		return model.Product{}, errors.New("description is required")
	}

	if product.Balance < 0 {
		return model.Product{}, errors.New("balance cannot be negative")
	}

	productId, err := pu.repository.CreateProductRepository(product)
	if err != nil {
		return model.Product{}, err
	}

	product.ID = productId
	return product, nil
}

func (pu *ProductUseCase) GetProductsUseCase() ([]model.Product, error) {
	return pu.repository.GetProductsRepository()
}

func (pu *ProductUseCase) GetProductByCodeUseCase(codeProduct string) (*model.Product, error) {
	product, err := pu.repository.GetProductByCodeRepository(codeProduct)
	if err != nil {
		return nil, err
	}

	return product, nil
}
