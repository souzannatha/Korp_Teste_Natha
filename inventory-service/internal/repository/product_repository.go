package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/souzannatha/Korp_Teste_Natha/inventory-service/internal/model"
)

type ProductRepository struct {
	connection *sql.DB
}

var ErrDuplicateCode = errors.New("product code already exists")

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) CreateProductRepository(ctx context.Context, product model.Product) (int, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	var id int
	query, err := pr.connection.PrepareContext(ctx, `
 INSERT INTO products (code, description, balance)
 VALUES ($1, $2, $3)
 RETURNING id
 `)
	if err != nil {
		return 0, err
	}
	defer query.Close()
	err = query.QueryRowContext(ctx, product.Code, product.Description, product.Balance).Scan(&id)
	if err != nil {
		var databaseErr *pq.Error
		if errors.As(err, &databaseErr) && databaseErr.Code == "23505" {
			return 0, ErrDuplicateCode
		}
		return 0, err
	}
	return id, nil
}

func (pr *ProductRepository) GetProductsRepository(ctx context.Context) ([]model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query := "SELECT id, code, description, balance FROM products ORDER BY id"
	rows, err := pr.connection.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	productList := []model.Product{}
	for rows.Next() {
		var productObj model.Product
		err = rows.Scan(&productObj.ID, &productObj.Code, &productObj.Description, &productObj.Balance)
		if err != nil {
			return nil, err
		}
		productList = append(productList, productObj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return productList, nil
}

func (pr *ProductRepository) GetProductByCodeRepository(ctx context.Context, codeProduct string) (*model.Product, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	query, err := pr.connection.PrepareContext(ctx, `
 SELECT id, code, description, balance FROM products WHERE code = $1
 `)
	if err != nil {
		return nil, err
	}
	defer query.Close()
	var product model.Product
	err = query.QueryRowContext(ctx, codeProduct).Scan(&product.ID, &product.Code, &product.Description, &product.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}
