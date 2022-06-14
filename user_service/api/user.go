package api

import "github.com/shopspring/decimal"

type CreateUserRequest struct {
	Name string `json:"name"`
}

type User struct {
	ID      int             `json:"ID"`
	Name    string          `json:"name"`
	Balance decimal.Decimal `json:"balance"`
}

type CreateUserResponse User
type GetUserResponse User

type DepositBalanceRequest struct {
	ID     int             `json:"ID"`
	Amount decimal.Decimal `json:"amount"`
}

type DepositBalanceResponse User
