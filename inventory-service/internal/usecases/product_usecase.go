package usecase

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/model"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/repository"
)

type ProductUseCase struct {
	repository repository.ProductRepository
}

var (
	ErrCodeRequired        = errors.New("code is required")
	ErrCodeTooLong         = errors.New("code must have at most 50 characters")
	ErrDescriptionRequired = errors.New("description is required")
	ErrBalanceRequired     = errors.New("balance is required")
	ErrInvalidBalance      = errors.New("balance must be between 0 and 2147483647")
)

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return ProductUseCase{
		repository: repo,
	}
}

func (pu *ProductUseCase) CreateProductUseCase(ctx context.Context, product model.Product) (model.Product, error) {
	product.Code = strings.TrimSpace(product.Code)
	product.Description = strings.TrimSpace(product.Description)
	if product.Code == "" {
		return model.Product{}, ErrCodeRequired
	}
	if utf8.RuneCountInString(product.Code) > 50 {
		return model.Product{}, ErrCodeTooLong
	}
	if product.Description == "" {
		return model.Product{}, ErrDescriptionRequired
	}
	if product.Balance < 0 || product.Balance > 2147483647 {
		return model.Product{}, ErrInvalidBalance
	}

	productId, err := pu.repository.CreateProductRepository(ctx, product)
	if err != nil {
		return model.Product{}, err
	}
	product.ID = productId
	return product, nil
}

func (pu *ProductUseCase) GetProductsUseCase(ctx context.Context) ([]model.Product, error) {
	return pu.repository.GetProductsRepository(ctx)
}

func (pu *ProductUseCase) GetProductByCodeUseCase(ctx context.Context, codeProduct string) (*model.Product, error) {
	codeProduct = strings.TrimSpace(codeProduct)
	if codeProduct == "" {
		return nil, ErrCodeRequired
	}
	if utf8.RuneCountInString(codeProduct) > 50 {
		return nil, ErrCodeTooLong
	}
	return pu.repository.GetProductByCodeRepository(ctx, codeProduct)
}
