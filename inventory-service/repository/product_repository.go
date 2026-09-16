package repository

import (
	"database/sql"
	"fmt"

	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) CreateProductRepository(product model.Product) (int, error) {
	var id int

	query, err := pr.connection.Prepare(`
    INSERT INTO products (code, description, balance)
    VALUES ($1, $2, $3)
    RETURNING id
`)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Code, product.Description, product.Balance).Scan(&id)

	if err != nil {

		fmt.Println(err)
		return 0, err
	}
	defer query.Close()

	return id, nil
}
