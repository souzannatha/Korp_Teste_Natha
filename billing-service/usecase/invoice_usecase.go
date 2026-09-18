package usecase

import (
	"errors"
	"strings"

	"github.com/souzannatha/Korp_Teste_Natha/billing-service/model"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/repository"
)

type InvoiceUseCase struct {
	repository repository.InvoiceRepository
}

func NewInvoiceUseCase(repo repository.InvoiceRepository) InvoiceUseCase {
	return InvoiceUseCase{
		repository: repo,
	}
}

func (iuc *InvoiceUseCase) CreateInvoiceUseCase(invoice model.Invoice) (model.Invoice, error) {
	if len(invoice.Items) == 0 {
		return model.Invoice{}, errors.New("invoice must have at least one item")
	}

	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.ProductCode = strings.TrimSpace(item.ProductCode)

		if item.ProductCode == "" {
			return model.Invoice{}, errors.New("product code is required")
		}

		if item.Quantity <= 0 {
			return model.Invoice{}, errors.New("quantity must be greater than zero")
		}
	}
	return iuc.repository.CreateInvoiceRepository(invoice)
}

func (iuc *InvoiceUseCase) GetAllInvoicesUseCase() ([]model.Invoice, error) {
	return iuc.repository.GetAllInvoicesRepository()
}
