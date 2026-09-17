package model

type InvoiceStatus string

const (
	StatusOpen   InvoiceStatus = "open"
	StatusClosed InvoiceStatus = "closed"
)

type Invoice struct {
	ID     int           `json:"id"`
	Number int           `json:"number"`
	Status InvoiceStatus `json:"status"`
	Items  []InvoiceItem `json:"items"`
}
