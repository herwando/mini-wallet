package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Withdrawal struct {
	ID            string          `json:"id" db:"id"`
	WithdrawnBy   string          `json:"withdrawn_by" db:"withdrawn_by"`
	Status        int             `json:"-" db:"status"`
	StatusMessage string          `json:"status" db:"-"`
	WithdrawnAt   time.Time       `json:"-" db:"withdrawn_at"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	ReferenceID   string          `json:"reference_id" db:"reference_id"`
}

type PayloadWithdrawal struct {
	ReferenceID string          `json:"reference_id"`
	Amount      decimal.Decimal `json:"amount"`
}
