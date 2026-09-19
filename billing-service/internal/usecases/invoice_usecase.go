package usecase

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/client"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/model"
	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/repository"
)

type InvoiceUseCase struct {
	repository      repository.InvoiceRepository
	inventoryClient client.InventoryClient
}

var (
	ErrInvoiceWithoutItems = errors.New("invoice must have at least one item")
	ErrProductCodeRequired = errors.New("product code is required")
	ErrProductCodeTooLong  = errors.New("product code must have at most 50 characters")
	ErrInvalidQuantity     = errors.New("quantity must be between 1 and 2147483647")
	ErrInvalidInvoiceId    = errors.New("invoice id must be greater than zero")
	ErrInvoiceClosed       = errors.New("only open invoices can be printed")
)

func NewInvoiceUseCase(repo repository.InvoiceRepository, inventoryClient client.InventoryClient) InvoiceUseCase {
	return InvoiceUseCase{
		repository:      repo,
		inventoryClient: inventoryClient,
	}
}

func (iuc *InvoiceUseCase) CreateInvoiceUseCase(ctx context.Context, invoice model.Invoice) (model.Invoice, error) {
	if len(invoice.Items) == 0 {
		return model.Invoice{}, ErrInvoiceWithoutItems
	}
	quantities := map[string]int{}
	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.ProductCode = strings.TrimSpace(item.ProductCode)
		if item.ProductCode == "" {
			return model.Invoice{}, ErrProductCodeRequired
		}
		if utf8.RuneCountInString(item.ProductCode) > 50 {
			return model.Invoice{}, ErrProductCodeTooLong
		}
		if item.Quantity <= 0 || item.Quantity > 2147483647-quantities[item.ProductCode] {
			return model.Invoice{}, ErrInvalidQuantity
		}
		quantities[item.ProductCode] += item.Quantity
	}
	checked := map[string]bool{}
	for _, item := range invoice.Items {
		if checked[item.ProductCode] {
			continue
		}
		if err := iuc.inventoryClient.GetProductByCode(ctx, item.ProductCode); err != nil {
			return model.Invoice{}, err
		}
		checked[item.ProductCode] = true
	}
	return iuc.repository.CreateInvoiceRepository(ctx, invoice)
}

func (iuc *InvoiceUseCase) GetAllInvoicesUseCase(ctx context.Context) ([]model.Invoice, error) {
	return iuc.repository.GetAllInvoicesRepository(ctx)
}

func (iuc *InvoiceUseCase) GetInvoiceByIdUseCase(ctx context.Context, invoiceId int) (model.Invoice, error) {
	if invoiceId <= 0 {
		return model.Invoice{}, ErrInvalidInvoiceId
	}
	return iuc.repository.GetInvoiceByIdRepository(ctx, invoiceId)
}

func (iuc *InvoiceUseCase) PrintInvoiceUseCase(ctx context.Context, invoiceId int) (model.Invoice, error) {
	if invoiceId <= 0 {
		return model.Invoice{}, ErrInvalidInvoiceId
	}
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	tx, err := iuc.repository.BeginInvoiceTransactionRepository(ctx)
	if err != nil {
		return model.Invoice{}, err
	}
	defer tx.Rollback()
	invoice, err := iuc.repository.GetInvoiceForPrintRepository(ctx, tx, invoiceId)
	if err != nil {
		return model.Invoice{}, err
	}
	if invoice.Status != model.StatusOpen {
		return model.Invoice{}, ErrInvoiceClosed
	}
	if len(invoice.Items) == 0 {
		return model.Invoice{}, ErrInvoiceWithoutItems
	}
	deduction := model.StockDeduction{InvoiceID: invoice.ID, Items: invoice.Items}
	if err := iuc.inventoryClient.DeductStock(ctx, deduction); err != nil {
		return model.Invoice{}, err
	}
	if err := iuc.repository.CloseInvoiceRepository(ctx, tx, invoice.ID); err != nil {
		return model.Invoice{}, err
	}
	if err := tx.Commit(); err != nil {
		return model.Invoice{}, err
	}
	invoice.Status = model.StatusClosed
	return invoice, nil
}
