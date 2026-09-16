package model

type Product struct {
	ID          int    `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     int    `json:"balance"`
}
