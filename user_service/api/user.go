package api

import "github.com/shopspring/decimal"

type CreateUserRequest struct {
	Name string `json:"name"`
}

type CreateUserResponse struct {
	ID      int             `json:"ID"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}
