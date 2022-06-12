package domain

import (
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	ID        int             `db:"id"`
	Name      string          `db:"name"`
	Balance   decimal.Decimal `db:"balance"`
	CratedAt  time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}
