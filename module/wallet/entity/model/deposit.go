package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Deposit struct {
	ID            string          `json:"id" db:"id"`
	DepositedBy   string          `json:"deposited_by" db:"deposited_by"`
	Status        int             `json:"-" db:"status"`
	StatusMessage string          `json:"status" db:"-"`
	DepositedAt   time.Time       `json:"-" db:"deposited_at"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	ReferenceID   string          `json:"reference_id" db:"reference_id"`
}

type PayloadDeposit struct {
	ReferenceID string          `json:"reference_id"`
	Amount      decimal.Decimal `json:"amount"`
}
