package repository

import (
	"database/sql"
	"fmt"

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

func (iv *InvoiceRepository) GetAllInvoicesRepository() ([]model.Invoice, error) {
	query := `
	SELECT 
	id, number, status 
	FROM invoices 
	ORDER BY number
	`
	rows, err := iv.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Invoice{}, err
	}
	defer rows.Close()

	invoices := []model.Invoice{}

	for rows.Next() {
		var invoice model.Invoice

		err = rows.Scan(&invoice.ID, &invoice.Number, &invoice.Status)
		if err != nil {
			fmt.Println(err)
			return []model.Invoice{}, err
		}
		invoices = append(invoices, invoice)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range invoices {
		itemRows, err := iv.connection.Query(`
		SELECT 
		product_code, quantity 
		FROM invoice_items 
		WHERE invoice_id = $1
		ORDER BY id
	`, invoices[i].ID)

		if err != nil {
			return nil, err
		}

		invoices[i].Items = []model.InvoiceItem{}

		for itemRows.Next() {
			var item model.InvoiceItem

			if err := itemRows.Scan(&item.ProductCode, &item.Quantity); err != nil {
				itemRows.Close()
				return nil, err
			}

			invoices[i].Items = append(invoices[i].Items, item)
		}

		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, err
		}

		itemRows.Close()
	}

	rows.Close()
	return invoices, nil
}
