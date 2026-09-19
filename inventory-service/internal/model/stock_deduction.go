package model

type StockDeduction struct {
	InvoiceID int                  `json:"invoice_id"`
	Items     []StockDeductionItem `json:"items"`
}
type StockDeductionItem struct {
	ProductCode string `json:"product_code"`
	Quantity    int    `json:"quantity"`
}
