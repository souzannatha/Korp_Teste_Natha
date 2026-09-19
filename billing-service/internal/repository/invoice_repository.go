package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/souzannatha/Korp_Teste_Natha/billing-service/internal/model"
)

type InvoiceRepository struct {
	connection *sql.DB
}

func NewInvoiceRepository(connection *sql.DB) InvoiceRepository {
	return InvoiceRepository{
		connection: connection,
	}
}

func (iv *InvoiceRepository) CreateInvoiceRepository(ctx context.Context, invoice model.Invoice) (model.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	tx, err := iv.connection.BeginTx(ctx, nil)
	if err != nil {
		return model.Invoice{}, err
	}
	defer tx.Rollback()
	err = tx.QueryRowContext(ctx, `
 INSERT INTO invoices DEFAULT VALUES RETURNING id, number, status
 `).Scan(&invoice.ID, &invoice.Number, &invoice.Status)
	if err != nil {
		return model.Invoice{}, err
	}
	for _, item := range invoice.Items {
		_, err = tx.ExecContext(ctx, `
  INSERT INTO invoice_items (invoice_id, product_code, quantity)
  VALUES ($1, $2, $3)
  `, invoice.ID, item.ProductCode, item.Quantity)
		if err != nil {
			return model.Invoice{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return model.Invoice{}, err
	}
	return invoice, nil
}

func (iv *InvoiceRepository) GetAllInvoicesRepository(ctx context.Context) ([]model.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := `SELECT id, number, status FROM invoices ORDER BY number`
	rows, err := iv.connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	invoices := []model.Invoice{}
	for rows.Next() {
		var invoice model.Invoice
		err = rows.Scan(&invoice.ID, &invoice.Number, &invoice.Status)
		if err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	for i := range invoices {
		items, err := iv.GetInvoiceItemsRepository(ctx, invoices[i].ID)
		if err != nil {
			return nil, err
		}
		invoices[i].Items = items
	}
	return invoices, nil
}

func (iv *InvoiceRepository) GetInvoiceByIdRepository(ctx context.Context, invoiceId int) (model.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var invoice model.Invoice
	query := `SELECT id, number, status FROM invoices WHERE id = $1`
	err := iv.connection.QueryRowContext(ctx, query, invoiceId).Scan(&invoice.ID, &invoice.Number, &invoice.Status)
	if err != nil {
		return model.Invoice{}, err
	}
	invoice.Items, err = iv.GetInvoiceItemsRepository(ctx, invoice.ID)
	if err != nil {
		return model.Invoice{}, err
	}
	return invoice, nil
}

func (iv *InvoiceRepository) GetInvoiceItemsRepository(ctx context.Context, invoiceId int) ([]model.InvoiceItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	itemRows, err := iv.connection.QueryContext(ctx, `
 SELECT product_code, quantity FROM invoice_items
 WHERE invoice_id = $1 ORDER BY id
 `, invoiceId)
	if err != nil {
		return nil, err
	}
	defer itemRows.Close()
	items := []model.InvoiceItem{}
	for itemRows.Next() {
		var item model.InvoiceItem
		if err := itemRows.Scan(&item.ProductCode, &item.Quantity); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := itemRows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (iv *InvoiceRepository) BeginInvoiceTransactionRepository(ctx context.Context) (*sql.Tx, error) {
	return iv.connection.BeginTx(ctx, nil)
}

func (iv *InvoiceRepository) GetInvoiceForPrintRepository(ctx context.Context, tx *sql.Tx, invoiceId int) (model.Invoice, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var invoice model.Invoice
	err := tx.QueryRowContext(ctx, `
 SELECT id, number, status FROM invoices WHERE id = $1 FOR UPDATE
 `, invoiceId).Scan(&invoice.ID, &invoice.Number, &invoice.Status)
	if err != nil {
		return model.Invoice{}, err
	}
	itemRows, err := tx.QueryContext(ctx, `
 SELECT product_code, quantity FROM invoice_items
 WHERE invoice_id = $1 ORDER BY id
 `, invoice.ID)
	if err != nil {
		return model.Invoice{}, err
	}
	defer itemRows.Close()
	invoice.Items = []model.InvoiceItem{}
	for itemRows.Next() {
		var item model.InvoiceItem
		if err := itemRows.Scan(&item.ProductCode, &item.Quantity); err != nil {
			return model.Invoice{}, err
		}
		invoice.Items = append(invoice.Items, item)
	}
	if err := itemRows.Err(); err != nil {
		return model.Invoice{}, err
	}
	return invoice, nil
}

func (iv *InvoiceRepository) CloseInvoiceRepository(ctx context.Context, tx *sql.Tx, invoiceId int) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err := tx.ExecContext(ctx, `UPDATE invoices SET status = 'closed' WHERE id = $1`, invoiceId)
	return err
}
