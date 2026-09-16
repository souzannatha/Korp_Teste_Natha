package usecase

import (
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
	productId, err := pu.repository.CreateProductRepository(product)

	if err != nil {
		return model.Product{}, err
	}
	product.ID = productId

	return product, nil
}
