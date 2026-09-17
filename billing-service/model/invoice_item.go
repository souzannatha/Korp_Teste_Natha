package model

type InvoiceItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}
