package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type Wallet struct {
	ID            string          `json:"id" db:"id"`
	OwnedBy       string          `json:"owned_by" db:"owned_by"`
	Status        int             `json:"-" db:"status"`
	StatusMessage string          `json:"status" db:"-"`
	EnabledAt     time.Time       `json:"-" db:"enabled_at"`
	Balance       decimal.Decimal `json:"balance" db:"balance"`
}

type PayloadDisable struct {
	IsDisabled bool `json:"is_disabled"`
}
