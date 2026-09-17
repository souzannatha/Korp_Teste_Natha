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

func (pr *ProductRepository) GetProductsRepository() ([]model.Product, error) {
	query := "SELECT id, code, description,balance FROM products"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Code,
			&productObj.Description,
			&productObj.Balance)

		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}
		productList = append(productList, productObj)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	defer rows.Close()
	return productList, nil
}

func (pr *ProductRepository) GetProductByCodeRepository(idProduct string) (*model.Product, error) {
	query, err := pr.connection.Prepare(`SELECT id, code, description, balance
	FROM products 
	WHERE code = $1
	`)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product model.Product

	err = query.QueryRow(idProduct).Scan(
		&product.ID,
		&product.Code,
		&product.Description,
		&product.Balance,
	)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	defer query.Close()
	return &product, nil
}
