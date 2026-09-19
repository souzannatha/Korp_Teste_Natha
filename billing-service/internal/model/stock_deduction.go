package model

type StockDeduction struct {
	InvoiceID int           `json:"invoice_id"`
	Items     []InvoiceItem `json:"items"`
}
