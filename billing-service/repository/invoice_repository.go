package repository

import (
	"database/sql"

	"github.com/souzannatha/Korp_Teste_Natha/billing-service/model"
)

type InvoiceRepository struct {
	connection *sql.DB
}

func NewInvoiceRepository(connection *sql.DB) InvoiceRepository {
	return InvoiceRepository{
		connection: connection,
	}
}

func (iv *InvoiceRepository) CreateInvoiceRepository(
	invoice model.Invoice,
) (model.Invoice, error) {
	tx, err := iv.connection.Begin()
	if err != nil {
		return model.Invoice{}, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
        INSERT INTO invoices DEFAULT VALUES
        RETURNING id, number, status
    `).Scan(
		&invoice.ID,
		&invoice.Number,
		&invoice.Status,
	)
	if err != nil {
		return model.Invoice{}, err
	}

	for _, item := range invoice.Items {
		_, err = tx.Exec(`
            INSERT INTO invoice_items (
                invoice_id, product_code, quantity
            )
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
