package model

type CreateProduct struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Balance     *int   `json:"balance"`
}
