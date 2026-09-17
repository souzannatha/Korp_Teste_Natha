package models

type InvoiceItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}
