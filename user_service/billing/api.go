package billing

import "github.com/shopspring/decimal"

type UserDepositRequest struct {
	ID     int             `json:"id"`
	Amount decimal.Decimal `json:"amount"`
}
