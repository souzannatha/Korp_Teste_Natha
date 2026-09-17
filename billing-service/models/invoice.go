package models

type InvoiceStatus string

const (
	StatusOpen   InvoiceStatus = "open"
	StatusClosed InvoiceStatus = "closed"
)

type Invoice struct {
	ID     int           `json:"id"`
	Number string        `json:"number"`
	Status InvoiceStatus `json:"status"`
	Items  []InvoiceItem `json:"items"`
}
