package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/model"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/repository"
)

var (
	ErrInvalidInvoiceId      = errors.New("invoice id must be greater than zero")
	ErrDeductionWithoutItems = errors.New("stock deduction must have at least one item")
	ErrInvalidQuantity       = errors.New("quantity must be between 1 and 2147483647")
)

type StockUseCase struct {
	repository repository.StockRepository
}

func NewStockUseCase(repo repository.StockRepository) StockUseCase {
	return StockUseCase{
		repository: repo,
	}
}

func (su *StockUseCase) DeductStockUseCase(ctx context.Context, deduction model.StockDeduction) error {
	deduction, fingerprint, err := normalizeDeduction(deduction)
	if err != nil {
		return err
	}
	return su.repository.DeductStockRepository(ctx, deduction, fingerprint)
}

func normalizeDeduction(deduction model.StockDeduction) (model.StockDeduction, string, error) {
	if deduction.InvoiceID <= 0 {
		return deduction, "", ErrInvalidInvoiceId
	}
	if len(deduction.Items) == 0 {
		return deduction, "", ErrDeductionWithoutItems
	}
	quantities := map[string]int{}
	for _, item := range deduction.Items {
		code := strings.TrimSpace(item.ProductCode)
		if code == "" {
			return deduction, "", ErrCodeRequired
		}
		if utf8.RuneCountInString(code) > 50 {
			return deduction, "", ErrCodeTooLong
		}
		if item.Quantity <= 0 || item.Quantity > 2147483647-quantities[code] {
			return deduction, "", ErrInvalidQuantity
		}
		quantities[code] += item.Quantity
	}
	codes := make([]string, 0, len(quantities))
	for code := range quantities {
		codes = append(codes, code)
	}
	sort.Strings(codes)
	deduction.Items = []model.StockDeductionItem{}
	for _, code := range codes {
		deduction.Items = append(deduction.Items, model.StockDeductionItem{ProductCode: code, Quantity: quantities[code]})
	}
	content, err := json.Marshal(deduction.Items)
	if err != nil {
		return deduction, "", err
	}
	fingerprint := sha256.Sum256(content)
	return deduction, hex.EncodeToString(fingerprint[:]), nil
}
