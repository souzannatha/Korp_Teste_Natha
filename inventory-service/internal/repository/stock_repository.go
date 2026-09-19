package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/model"
)

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrInsufficientBalance = errors.New("insufficient stock balance")
	ErrDeductionConflict   = errors.New("invoice already has a different stock deduction")
)

type StockRepository struct {
	connection *sql.DB
}

func NewStockRepository(connection *sql.DB) StockRepository {
	return StockRepository{
		connection: connection,
	}
}

func (sr *StockRepository) DeductStockRepository(ctx context.Context, deduction model.StockDeduction, fingerprint string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := sr.connection.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `
 INSERT INTO stock_deductions (invoice_id, fingerprint)
 VALUES ($1, $2)
 ON CONFLICT (invoice_id) DO NOTHING
 `, deduction.InvoiceID, fingerprint)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		var previousFingerprint string
		err = tx.QueryRowContext(ctx, `SELECT fingerprint FROM stock_deductions WHERE invoice_id = $1`, deduction.InvoiceID).Scan(&previousFingerprint)
		if err != nil {
			return err
		}
		if previousFingerprint != fingerprint {
			return ErrDeductionConflict
		}
		return tx.Commit()
	}

	for _, item := range deduction.Items {
		result, err = tx.ExecContext(ctx, `
  UPDATE products SET balance = balance - $1
  WHERE code = $2 AND balance >= $1
  `, item.Quantity, item.ProductCode)
		if err != nil {
			return err
		}
		affected, err = result.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			var exists bool
			err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM products WHERE code = $1)`, item.ProductCode).Scan(&exists)
			if err != nil {
				return err
			}
			if !exists {
				return ErrProductNotFound
			}
			return ErrInsufficientBalance
		}
	}
	return tx.Commit()
}
